//+build mage

package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	K8S_VERSION = "v1.20.0"
)

var (
	kubectl = sh.OutCmd("./build/tmp/kubectl", "--kubeconfig", "./build/tmp/kubeconfig")
	kind    = sh.OutCmd("./build/tmp/kind")
)

// Run executes the project with default settings
func Run() {
	sh.RunV("go", "run", "./cmd")
}

// Build generates a mimic binary in the root of the repository
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "mimic", "./cmd")
}

// Clean removes any remnants of the build/test process
func Clean() error {
	mg.Deps(ensureKind)
	if err := sh.Rm("mimic"); err != nil {
		return err
	}
	kind("delete", "cluster", "--name", fmt.Sprintf("mimic-%s", K8S_VERSION))
	if err := sh.Rm("./build"); err != nil {
		return err
	}
	return nil
}

// Lint runs golangci-lint against the project
func Lint() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return sh.RunV("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", pwd), "-w", "/app", "golangci/golangci-lint:v1.39.0", "golangci-lint", "run")
}

// Docker generates a docker image with a "latest" tag
func Docker() error {
	return DockerTag("latest")
}

// DockerTag generates a docker tag with a specific tag
func DockerTag(tag string) error {
	fmt.Println("[Image Build Tag] Starting")
	if err := sh.RunV("docker", "build", "-t", fmt.Sprintf("mimic:%v", tag), "."); err != nil {
		return err
	}
	fmt.Println("[Image Build Tag] Complete")
	return nil
}

// Deploy deploys all of the required Kubernetes manifests for Mimic into a KiND cluster.
func Deploy() error {
	mg.Deps(ensureKubectl, ensureCluster, ensureNamespace, generateCerts)

	fmt.Println("[Deploy] Deploying Kubernetes resources")
	if _, err := kubectl("apply", "-f", "./deploy/manifests"); err != nil {
		return err
	}

	mg.Deps(Update)

	fmt.Println("[Deploy] Complete")
	return nil
}

// Update deploys a newly built docker image of Mimic into a KiND cluster.
func Update() error {
	fmt.Println("[Update] Starting")
	tag := fmt.Sprint(time.Now().Unix())
	mg.Deps(ensureCluster, ensureKubectl, mg.F(pushImage, tag))

	fmt.Printf("[Update] Patching deployment to mimic:%s\n", tag)
	if _, err := kubectl("-n", "mimic", "set", "image", "deployment/mimic", fmt.Sprintf("mimic=mimic:%s", tag)); err != nil {
		return err
	}

	fmt.Println("[Minikube Update] Complete")
	return nil
}

// DeployHarbor will deploy the Harbor Helm Chart into the KiND cluster for testing the Harbor integration.
func DeployHarbor() error {
	mg.Deps(ensureCluster, ensureHelm)
	fmt.Println("[Harbor] Starting")
	envs := map[string]string{
		"KUBECONFIG":       "./build/tmp/kubeconfig",
		"HELM_CACHE_HOME":  "./build/tmp/helm_home/cache",
		"HELM_CONFIG_HOME": "./build/tmp/helm_home/config",
		"HELM_DATA_HOME":   "./build/tmp/helm_home/data",
		"HELM_KUBECONTEXT": fmt.Sprintf("kind-mimic-%s", K8S_VERSION),
	}
	if err := sh.RunWithV(envs, "./build/tmp/helm", "repo", "add", "harbor", "https://helm.goharbor.io"); err != nil {
		return err
	}
	if err := sh.RunWithV(envs, "./build/tmp/helm", "upgrade", "--install", "--create-namespace", "-n", "harbor",
		"--set", "expose.tls.enabled=false",
		"--set", "expose.type=clusterIP",
		"--set", "expose.clusterIP.name=localhost",
		"--set", "externalURL=http://localhost:8080",
		"mimic-harbor", "harbor/harbor"); err != nil {
		return err
	}
	fmt.Println("[Harbor] Complete")
	return nil
}

// GenerateHarborSwagger uses go-swagger to generate Harbor API bindings.
func GenerateHarborSwagger() error {
	mg.Deps(ensureDocker)

	if err := os.MkdirAll("./build/tmp", fs.FileMode(0775)); err != nil {
		return err
	}

	if err := os.MkdirAll("./internal/harbor", 0775); err != nil {
		return err
	}

	if err := downloadFile("./build/tmp/swagger.yaml", "https://raw.githubusercontent.com/goharbor/harbor/v2.2.1/api/v2.0/swagger.yaml"); err != nil {
		return err
	}

	if err := sh.Rm("./internal/harbor/models"); err != nil {
		return err
	}

	if err := sh.Rm("./internal/harbor/client"); err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := sh.RunV("docker", "run", "--rm", "-it",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", fmt.Sprintf("%s:%s", pwd, pwd),
		"-w", fmt.Sprintf("%s/build", pwd),
		"quay.io/goswagger/swagger",
		"generate", "client",
		"--target", "../internal/harbor",
		"--spec", "./tmp/swagger.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

func generateCerts() {
	mg.Deps(ensureCluster, ensureNamespace, ensureKubectl)
	fmt.Println("[Generate Certs] Starting")
	fmt.Println("[Generate Certs] Ensuring certificates are present on cluster")
	envs := map[string]string{
		"KUBECONFIG": "./build/tmp/kubeconfig",
	}
	sh.RunWith(envs, "./deploy/scripts/webhook-create-signed-cert.sh", "--service", "mimic", "--secret", "mimic-certs", "--namespace", "mimic")
	sh.RunWith(envs, "./deploy/scripts/webhook-patch-ca-bundle.sh", "./deploy/manifests/templates/mutatingwebhookconfiguration.yaml", "./deploy/manifests/mutatingwebhookconfiguration-cabundle.yaml")
	fmt.Println("[Generate Certs] Complete")
}

func pushImage(tag string) {
	mg.Deps(ensureKind, ensureCluster, mg.F(DockerTag, tag))
	kind("load", "docker-image", "--name", fmt.Sprintf("mimic-%s", K8S_VERSION), fmt.Sprintf("mimic:%s", tag))
}

func ensureDocker() error {
	err := sh.Run("docker", "info")
	if sh.ExitStatus(err) > 0 {
		return fmt.Errorf("Docker appears to not be running, please ensure docker is installed and running")
	}
	fmt.Println("Docker appears to be running")
	return nil
}

func ensureKubectl() error {
	if _, err := os.Stat("./build/tmp/kubectl"); os.IsNotExist(err) {
		fmt.Println("kubectl not present in build directory, fetching")
		if err := os.MkdirAll("./build/tmp", fs.FileMode(0775)); err != nil {
			return err
		}
		if err := downloadFile("./build/tmp/kubectl", fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", K8S_VERSION)); err != nil {
			return err
		}
		os.Chmod("./build/tmp/kubectl", fs.FileMode(0755))
		fmt.Println("kubectl fetched")
	}
	return nil
}

func ensureHelm() error {
	if _, err := os.Stat("./build/tmp/helm"); os.IsNotExist(err) {
		fmt.Println("Helm not present in build directory, fetching")
		if err := os.MkdirAll("./build/tmp", fs.FileMode(0775)); err != nil {
			return err
		}
		if err := downloadFile("./build/tmp/helm.tgz", "https://get.helm.sh/helm-v3.5.4-linux-amd64.tar.gz"); err != nil {
			return err
		}
		if err := sh.Run("tar", "xvf", "./build/tmp/helm.tgz", "--strip-components=1", "-C", "./build/tmp", "linux-amd64/helm"); err != nil {
			return err
		}
		fmt.Println("helm fetched")
	}
	return nil
}

func ensureKind() error {
	mg.Deps(ensureDocker)
	if _, err := os.Stat("./build/tmp/kind"); os.IsNotExist(err) {
		fmt.Println("KiND not present in build directory, fetching")
		if err := os.MkdirAll("./build/tmp", fs.FileMode(0775)); err != nil {
			return err
		}
		if err := downloadFile("./build/tmp/kind", "https://kind.sigs.k8s.io/dl/v0.10.0/kind-linux-amd64"); err != nil {
			return err
		}
		os.Chmod("./build/tmp/kind", fs.FileMode(0755))
		fmt.Println("KiND fetched")
	}
	return nil
}

func ensureCluster() error {
	mg.Deps(ensureKind, ensureDocker)
	output, err := kind("get", "clusters")
	if err != nil {
		return err
	}
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	clusterRunning := false
	for _, line := range lines {
		if line == fmt.Sprintf("mimic-%s", K8S_VERSION) {
			fmt.Println("KiND cluster running")
			clusterRunning = true
		}
	}
	if !clusterRunning {
		fmt.Println("Starting KiND cluster")
		kind("create", "cluster", "--name", fmt.Sprintf("mimic-%s", K8S_VERSION), "--image", fmt.Sprintf("kindest/node:%s", K8S_VERSION))
	}

	fmt.Println("Configuring kubeconfig for KiND cluster at ./build/tmp/kubeconfig")
	kubeconfig, err := kind("get", "kubeconfig", "--name", fmt.Sprintf("mimic-%s", K8S_VERSION))
	if err != nil {
		return err
	}

	file, err := os.Create("./build/tmp/kubeconfig")
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Chmod(fs.FileMode(0600)); err != nil {
		return err
	}
	if _, err := file.WriteString(kubeconfig); err != nil {
		return err
	}
	return nil
}

func ensureNamespace() error {
	mg.Deps(ensureKubectl, ensureCluster)
	_, err := kubectl("apply", "-f", "./deploy/manifests/namespace.yaml")
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
