package kexec

import (
	"fmt"
	"os/exec"
	"strings"
)

// Inyectables para test
var lookPath = exec.LookPath
var command = exec.Command

// CheckKubectlAvailable verifies that the 'kubectl' command is available in the system PATH.
func CheckKubectlAvailable() error {
	// This check ensures that all kubectl-based operations can proceed.
	if _, err := lookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl not found in PATH. Please install kubectl and ensure it is available in your PATH")
	}
	return nil
}

// CheckKubeCluster verifies that the user can connect to the Kubernetes cluster using kubectl.
func CheckKubeCluster() error {
	cmd := command("kubectl", "cluster-info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to connect to Kubernetes cluster: %v", err)
	}
	return nil
}

// CheckPodRunning verifies that the specified pod exists and is in the 'Running' phase.
func CheckPodRunning(namespace, pod string) error {
	cmd := command("kubectl", "get", "pod", "-n", namespace, pod, "-o", "jsonpath={.status.phase}")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get pod status: %v", err)
	}
	status := strings.TrimSpace(string(out))
	if status != "Running" {
		return fmt.Errorf("pod %s in namespace %s is not Running (status: %s)", pod, namespace, status)
	}
	return nil
}

// CheckExecPermission verifies that the user can execute commands in the specified pod/container.
// Returns an error if exec is not permitted (e.g., due to RBAC or ServiceAccount restrictions).
func CheckExecPermission(namespace, pod, container string) error {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--", "echo", "test")
	cmd := command("kubectl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl exec failed (check your RBAC/permissions): %v", err)
	}
	return nil
}

// CheckWritable verifies that the specified path in the pod is writable by attempting to create and remove a test file.
func CheckWritable(namespace, pod, container, path string) error {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	testFile := fmt.Sprintf("%s/.netdbg_write_test", path)
	args = append(args, "--", "sh", "-c", fmt.Sprintf("touch %s && rm %s", testFile, testFile))
	cmd := command("kubectl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot write to %s in pod: %v", path, err)
	}
	return nil
}

// CheckBinNotExists verifies that the netdbg binary does not already exist at the target path in the pod.
func CheckBinNotExists(namespace, pod, container, binPath string) error {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--", "test", "!", "-f", binPath)
	cmd := command("kubectl", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("netdbg binary already exists at %s in pod", binPath)
	}
	return nil
}

// CheckGoAvailable verifies that the 'go' command is available in the system PATH.
// Returns an error if Go is not found (required for --mode build).
func CheckGoAvailable() error {
	if _, err := lookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed or not in PATH, required for --mode build")
	}
	return nil
}

// GetPodArchOS retrieves the architecture and OS of the pod by running 'uname -m' and 'uname -s'.
func GetPodArchOS(namespace, pod, container string) (arch, osname string, err error) {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	argsArch := append(args, "--", "uname", "-m")
	cmdArch := command("kubectl", argsArch...)
	outArch, err := cmdArch.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get pod architecture: %v", err)
	}
	argsOS := append(args, "--", "uname", "-s")
	cmdOS := command("kubectl", argsOS...)
	outOS, err := cmdOS.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get pod OS: %v", err)
	}
	return strings.TrimSpace(string(outArch)), strings.TrimSpace(string(outOS)), nil
}
