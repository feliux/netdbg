package kexec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/feliux/netdbg/internal/logger"
)

// Inyectable for testing: allows overriding exec.Command in tests
var kubectlCommand = exec.Command

// KubectlCp copies a file from localBin to the pod at binPath.
func KubectlCp(namespace, pod, container, localBin, binPath string, dryRun bool) error {
	args := []string{"cp", localBin, fmt.Sprintf("%s/%s:%s", namespace, pod, binPath)}
	if container != "" {
		args = append(args, "-c", container)
	}
	logger.Info("kubectl cp", "args", args)
	if dryRun {
		logger.Info("[DRY-RUN] would copy netdbg binary to pod", "kubectl_args", strings.Join(args, " "))
		return nil
	}
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("kubectl cp failed", "error", err)
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
	logger.Info("kubectl exec chmod", "args", args)
	if dryRun {
		logger.Info("[DRY-RUN] would set executable permissions", "kubectl_args", strings.Join(args, " "))
		return nil
	}
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("kubectl exec chmod failed", "error", err)
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
	logger.Info("kubectl exec netdbg", "args", cmdArgs)
	if dryRun {
		logger.Info("[DRY-RUN] would execute in pod", "kubectl_args", strings.Join(cmdArgs, " "))
		return nil
	}
	cmd := kubectlCommand("kubectl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		logger.Error("kubectl exec netdbg failed", "error", err)
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
	logger.Info("kubectl exec cleanup", "args", args)
	if dryRun {
		logger.Info("[DRY-RUN] would cleanup binary from pod", "kubectl_args", strings.Join(args, " "))
		return nil
	}
	cmd := kubectlCommand("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("kubectl exec cleanup failed", "error", err)
		return err
	}
	return nil
}
