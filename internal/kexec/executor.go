package kexec

import (
	"context"
	"fmt"
	"os"

	"github.com/feliux/netdbg/internal/logger"
)

// Executor defines the interface for executing kexec operations.
type Executor interface {
	Execute(ctx context.Context, opts *Options, args []string) error
}

// DefaultExecutor implements the Executor interface with the standard kexec logic.
type DefaultExecutor struct{}

// Execute performs the kexec operation: verifies environment, prepares the binary,
// uploads it to the pod, executes the command, and handles cleanup as needed.
func (e *DefaultExecutor) Execute(ctx context.Context, opts *Options, args []string) error {
	localBin := "./netdbg"
	binSource := fmt.Sprintf("%s/%s/netdbg-%s-%s", opts.ReleaseBaseURL, opts.Version, opts.Goos, opts.Goarch)

	// 1. Check that kubectl is available
	logger.Info("checking kubectl availability")
	if !opts.DryRun {
		if err := CheckKubectlAvailable(); err != nil {
			logger.Debug("kubectl not available", "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check for kubectl in PATH")
	}

	// 2. Check cluster connectivity
	logger.Info("checking Kubernetes cluster connectivity")
	if !opts.DryRun {
		if err := CheckKubeCluster(); err != nil {
			logger.Debug("Kubernetes cluster not reachable", "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check Kubernetes cluster connectivity")
	}

	// 3. Check that the pod exists and is running
	logger.Info("checking pod status", "namespace", opts.Namespace, "pod", opts.Pod)
	if !opts.DryRun {
		if err := CheckPodRunning(opts.Namespace, opts.Pod); err != nil {
			logger.Debug("pod not running", "namespace", opts.Namespace, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check that pod is running", "namespace", opts.Namespace, "pod", opts.Pod)
	}

	// 4. Check exec permission
	logger.Info("checking exec permission", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container)
	if !opts.DryRun {
		if err := CheckExecPermission(opts.Namespace, opts.Pod, opts.Container); err != nil {
			logger.Debug("kubectl exec not permitted", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check exec permission on pod", "pod", opts.Pod)
	}

	// 5. Check write permission in the pod
	logger.Info("checking write permission", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "path", opts.BinPath)
	if !opts.DryRun {
		if err := CheckWritable(opts.Namespace, opts.Pod, opts.Container, opts.BinPath); err != nil {
			logger.Debug("path not writable in pod", "path", opts.BinPath, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check write permission at path in pod", "path", opts.BinPath, "pod", opts.Pod)
	}

	// 6. Check if binary already exists in the pod
	logger.Info("checking binary presence in pod", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "path", opts.BinPath)
	if !opts.DryRun {
		if err := CheckBinNotExists(opts.Namespace, opts.Pod, opts.Container, opts.BinPath); err != nil {
			logger.Debug("netdbg binary already exists in pod", "path", opts.BinPath, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Debug("[dry-run] would check if binary exists at path in pod", "path", opts.BinPath, "pod", opts.Pod)
	}

	// 7. If mode is build, check Go is available
	if opts.Mode == "build" {
		logger.Info("checking Go availability for build")
	}
	if opts.Mode == "build" && !opts.DryRun {
		if err := CheckGoAvailable(); err != nil {
			logger.Debug("Go not available for build", "error", err)
			return err
		}
	} else if opts.Mode == "build" && opts.DryRun {
		logger.Debug("[dry-run] would check if Go is available for local build")
	}

	// 8. Detect pod architecture/OS
	logger.Info("detecting pod architecture and OS", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container)
	if !opts.DryRun {
		arch, osname, err := GetPodArchOS(opts.Namespace, opts.Pod, opts.Container)
		if err != nil {
			logger.Debug("could not determine pod architecture/OS", "error", err)
		} else {
			logger.Debug("pod architecture and OS", "arch", arch, "os", osname)
			// Optionally compare with opts.Goarch/Goos and warn if mismatch
		}
	} else {
		logger.Debug("[dry-run] would detect pod architecture and OS")
	}

	// 9. Prepare the netdbg binary (build or download)
	logger.Info("preparing netdbg binary", "mode", opts.Mode, "goos", opts.Goos, "goarch", opts.Goarch, "version", opts.Version)
	if opts.Mode == "download" {
		if err := DownloadNetdbgBinary(binSource, localBin, opts.DryRun); err != nil {
			return err
		}
	} else if opts.Mode == "build" {
		if err := BuildNetdbgBinary(opts.Goos, opts.Goarch, localBin, opts.DryRun); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown mode: %s", opts.Mode)
	}

	// 10. Copy the binary to the pod
	logger.Info("copying netdbg binary to pod", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "dest", opts.BinPath)
	if err := KubectlCp(opts.Namespace, opts.Pod, opts.Container, localBin, opts.BinPath, opts.DryRun); err != nil {
		return err
	}

	// 11. Set executable permissions
	logger.Info("setting executable permissions in pod", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "path", opts.BinPath)
	if err := KubectlChmod(opts.Namespace, opts.Pod, opts.Container, opts.BinPath, opts.DryRun); err != nil {
		return err
	}

	// 12. Execute netdbg in the pod
	logger.Info("executing netdbg in pod", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container)
	if err := KubectlExec(opts.Namespace, opts.Pod, opts.Container, opts.BinPath, args, opts.DryRun); err != nil {
		return err
	}

	// 13. Cleanup if requested
	if opts.Cleanup {
		logger.Info("removing netdbg binary from pod", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "path", opts.BinPath)
		if err := KubectlCleanup(opts.Namespace, opts.Pod, opts.Container, opts.BinPath, opts.DryRun); err != nil {
			return err
		}
	}

	// 14. Optionally cleanup local binary
	if !opts.KeepLocal && !opts.DryRun {
		if _, err := os.Stat(localBin); err == nil {
			_ = os.Remove(localBin)
		}
	}

	return nil
}
