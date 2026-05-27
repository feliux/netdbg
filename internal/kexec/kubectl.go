package kexec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/feliux/netdbg/internal/logger"
)

// Injectable for testing: allows overriding exec.Command in tests
var kubectlCommand = exec.Command

// KubectlCp copies a file from localBin to the pod at binPath.
func KubectlCp(namespace, pod, container, localBin, binPath string, dryRun bool) error {
	args := []string{"cp", localBin, fmt.Sprintf("%s/%s:%s", namespace, pod, binPath)}
	if container != "" {
		args = append(args, "-c", container)
	}
	if dryRun {
		logger.Debug("[dry-run] copying netdbg binary to pod", "namespace", namespace, "pod", pod, "container", container, "src", localBin, "dest", binPath, "kubectl_args", strings.Join(args, " "))
		return nil
	}
	logger.Debug("copying netdbg binary to pod", "namespace", namespace, "pod", pod, "container", container, "src", localBin, "dest", binPath)
	logger.Debug("kubectl cp args", "args", strings.Join(args, " "))
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Debug("kubectl cp failed", "error", err)
		return err
	}
	return nil
}

// KubectlChmod sets executable permissions on the binary in the pod.
func KubectlChmod(namespace, pod, container, binPath string, dryRun bool) error {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--", "chmod", "+x", binPath)
	if dryRun {
		logger.Debug("[dry-run] setting executable permissions in pod", "namespace", namespace, "pod", pod, "container", container, "path", binPath, "kubectl_args", strings.Join(args, " "))
		return nil
	}
	logger.Debug("setting executable permissions in pod", "namespace", namespace, "pod", pod, "container", container, "path", binPath)
	logger.Debug("kubectl exec chmod args", "args", strings.Join(args, " "))
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Debug("kubectl exec chmod failed", "error", err)
		return err
	}
	return nil
}

// KubectlExec runs the netdbg command in the pod.
func KubectlExec(namespace, pod, container, binPath string, args []string, dryRun bool) error {
	cmdArgs := []string{"exec", "-n", namespace, pod}
	if container != "" {
		cmdArgs = append(cmdArgs, "-c", container)
	}
	cmdArgs = append(cmdArgs, "--", binPath)
	cmdArgs = append(cmdArgs, args...)
	if dryRun {
		logger.Debug("[dry-run] executing netdbg in pod", "namespace", namespace, "pod", pod, "container", container, "bin_path", binPath, "args", args, "kubectl_args", strings.Join(cmdArgs, " "))
		return nil
	}
	logger.Debug("executing netdbg in pod", "namespace", namespace, "pod", pod, "container", container, "bin_path", binPath, "args", args)
	logger.Debug("kubectl exec netdbg args", "args", strings.Join(cmdArgs, " "))
	cmd := kubectlCommand("kubectl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		logger.Debug("kubectl exec netdbg failed", "error", err)
		return err
	}
	return nil
}

// KubectlCleanup removes the binary from the pod.
func KubectlCleanup(namespace, pod, container, binPath string, dryRun bool) error {
	args := []string{"exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--", "rm", "-f", binPath)
	if dryRun {
		logger.Debug("[dry-run] removing netdbg binary from pod", "namespace", namespace, "pod", pod, "container", container, "path", binPath, "kubectl_args", strings.Join(args, " "))
		return nil
	}
	logger.Debug("removing netdbg binary from pod", "namespace", namespace, "pod", pod, "container", container, "path", binPath)
	logger.Debug("kubectl exec cleanup args", "args", strings.Join(args, " "))
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Debug("kubectl exec cleanup failed", "error", err)
		return err
	}
	return nil
}
