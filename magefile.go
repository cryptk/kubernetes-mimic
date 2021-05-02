//+build mage

package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Run() {
	sh.RunV("go", "run", "./cmd")
}

func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "mimic", "./cmd")
}

func Clean() error {
	if err := sh.Rm("mimic"); err != nil {
		return err
	}
	if err := sh.Rm("./build"); err != nil {
		return err
	}
	return nil
}

func ImageBuild() error {
	return sh.RunV("docker", "build", "-t", "mimic", ".")
}

func ImageBuildTag(tag string) error {
	fmt.Println("[Image Build Tag] Starting")
	if err := sh.RunV("docker", "build", "-t", fmt.Sprintf("mimic:%v", tag), "."); err != nil {
		return err
	}
	fmt.Println("[Image Build Tag] Complete")
	return nil
}

func ensureMinikube() error {
	fmt.Println("[Ensure Minikube] Starting")
	err := sh.Run("minikube", "status")
	code := sh.ExitStatus(err)
	switch code {
	case 0: // Minikube already running
		break
	case 7, 85: // 7 == minikube stopped, 85 == minikube not created
		sh.Run("minikube", "start")
	default:
		fmt.Printf("Unknown `minikube status` return code: %v\n", code)
		return err
	}
	sh.Run("minikube", "update-context")
	fmt.Println("[Ensure Minikube] Complete")
	return nil
}

func generateCerts() {
	fmt.Println("[Generate Certs] Starting")
	fmt.Println("[Generate Certs] Ensuring certificates are present on cluster")
	sh.Run("./deploy/scripts/webhook-create-signed-cert.sh", "--service", "mimic", "--secret", "mimic-certs", "--namespace", "mimic")
	sh.Run("./deploy/scripts/webhook-patch-ca-bundle.sh", "./deploy/manifests/templates/mutatingwebhookconfiguration.yaml", "./deploy/mutatingwebhookconfiguration-cabundle.yaml")
	fmt.Println("[Generate Certs] Complete")
}

func Lint() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return sh.RunV("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", pwd), "-w", "/app", "golangci/golangci-lint:v1.39.0", "golangci-lint", "run")
}

// Deploy Mimic into a Minikube cluster.  Assumes that
func MKDeploy() error {
	fmt.Println("[Minikube Deploy] Starting")
	mg.Deps(ensureMinikube)

	fmt.Println("[Minikube Deploy] Deploying Kubernetes resources")
	if err := sh.Run("minikube", "kubectl", "--", "apply", "-f", "./deploy/manifests/namespace.yaml"); err != nil {
		return err
	}
	mg.Deps(generateCerts)
	sh.Run("minikube", "kubectl", "--", "apply", "-f", "./deploy/manifests")

	fmt.Println("[Minikube Deploy] Complete")
	return nil
}

// Update Mimic into a Minikube cluster.  Assumes that
func MKUpdate() error {
	fmt.Println("[Minikube Update] Starting")
	tag := fmt.Sprint(time.Now().Unix())
	mg.Deps(ensureMinikube, mg.F(ImageBuildTag, tag))

	fmt.Printf("[Minikube Update] Sending mimic:%s to minikube\n", tag)
	if err := sh.Run("minikube", "image", "load", fmt.Sprintf("mimic:%s", tag)); err != nil {
		return err
	}
	fmt.Printf("[Minikube Update] Patching deployment to mimic:%s\n", tag)
	sh.Run("minikube", "kubectl", "--", "-n", "mimic", "set", "image", "deployment/mimic", fmt.Sprintf("mimic=mimic:%s", tag))

	fmt.Println("[Minikube Update] Complete")
	return nil
}

func MKHarbor() error {
	mg.Deps(ensureMinikube)
	fmt.Println("[Minikube Harbor] Starting")
	if err := sh.Run("helm", "repo", "add", "harbor", "https://helm.goharbor.io"); err != nil {
		return err
	}
	if err := sh.Run("helm", "upgrade", "--install", "--create-namespace", "-n", "mimic-harbor",
		"--set", "expose.tls.enabled=false",
		"--set", "expose.type=clusterIP",
		"--set", "expose.clusterIP.name=localhost",
		"--set", "externalURL=http://localhost:8080",
		"mimic-harbor", "harbor/harbor"); err != nil {
		return err
	}
	fmt.Println("[Minikube Harbor] Complete")
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

// GenerateHarborSwagger uses go-swagger to generate Harbor API bindings.
func GenerateHarborSwagger() error {
	// mg.Deps(ensureMinikube)

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
