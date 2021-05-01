//+build mage

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Run() {
	sh.RunV("go", "run", "./cmd")
}

// This builds the project.
// This is a multi-line description.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "mimic", "./cmd")
}

func Clean() error {
	return sh.Rm("mimic")
}

func ImageBuild() error {
	return sh.Run("docker", "build", "-t", "mimic", ".")
}

func ImageBuildTag(tag string) error {
	fmt.Println("[Image Build Tag] Starting")
	if err := sh.Run("docker", "build", "-t", fmt.Sprintf("mimic:%v", tag), "."); err != nil {
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
	sh.Run("./deploy/webhook-create-signed-cert.sh", "--service", "mimic", "--secret", "mimic-certs", "--namespace", "mimic")
	sh.Run("./deploy/webhook-patch-ca-bundle.sh", "./deploy/templates/mutatingwebhookconfiguration.yaml", "./deploy/mutatingwebhookconfiguration-cabundle.yaml")
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
	if err := sh.Run("minikube", "kubectl", "--", "apply", "-f", "./deploy/namespace.yaml"); err != nil {
		return err
	}
	mg.Deps(generateCerts)
	sh.Run("minikube", "kubectl", "--", "apply", "-f", "./deploy")

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
