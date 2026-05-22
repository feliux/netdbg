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
	if !opts.DryRun {
		if err := CheckKubectlAvailable(); err != nil {
			logger.Error("kubectl not available", "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check for kubectl in PATH")
	}

	// 2. Check cluster connectivity
	if !opts.DryRun {
		if err := CheckKubeCluster(); err != nil {
			logger.Error("kubernetes cluster not reachable", "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check Kubernetes cluster connectivity")
	}

	// 3. Check that the pod exists and is running
	if !opts.DryRun {
		if err := CheckPodRunning(opts.Namespace, opts.Pod); err != nil {
			logger.Error("pod not running", "namespace", opts.Namespace, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check that pod is running", "namespace", opts.Namespace, "pod", opts.Pod)
	}

	// 4. Check exec permission
	if !opts.DryRun {
		if err := CheckExecPermission(opts.Namespace, opts.Pod, opts.Container); err != nil {
			logger.Error("kubectl exec not permitted", "namespace", opts.Namespace, "pod", opts.Pod, "container", opts.Container, "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check exec permission on pod", "pod", opts.Pod)
	}

	// 5. Check write permission in the pod
	if !opts.DryRun {
		if err := CheckWritable(opts.Namespace, opts.Pod, opts.Container, opts.BinPath); err != nil {
			logger.Error("path not writable in pod", "path", opts.BinPath, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check write permission at path in pod", "path", opts.BinPath, "pod", opts.Pod)
	}

	// 6. Check if binary already exists in the pod
	if !opts.DryRun {
		if err := CheckBinNotExists(opts.Namespace, opts.Pod, opts.Container, opts.BinPath); err != nil {
			logger.Error("netdbg binary already exists in pod", "path", opts.BinPath, "pod", opts.Pod, "error", err)
			return err
		}
	} else {
		logger.Info("[DRY-RUN] would check if binary exists at path in pod", "path", opts.BinPath, "pod", opts.Pod)
	}

	// 7. If mode is build, check Go is available
	if opts.Mode == "build" && !opts.DryRun {
		if err := CheckGoAvailable(); err != nil {
			logger.Error("golang not available for build", "error", err)
			return err
		}
	} else if opts.Mode == "build" && opts.DryRun {
		logger.Info("[DRY-RUN] would check if Go is available for local build")
	}

	// 8. Detect pod architecture/OS
	if !opts.DryRun {
		arch, osname, err := GetPodArchOS(opts.Namespace, opts.Pod, opts.Container)
		if err != nil {
			logger.Error("could not determine pod architecture/OS", "error", err)
		} else {
			logger.Info("pod architecture and OS", "arch", arch, "os", osname)
			// Optionally compare with opts.Goarch/Goos and warn if mismatch
		}
	} else {
		logger.Info("[DRY-RUN] would detect pod architecture and OS")
	}

	// 9. Prepare the netdbg binary (build or download)
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
	if err := KubectlCp(opts.Namespace, opts.Pod, opts.Container, localBin, opts.BinPath, opts.DryRun); err != nil {
		return err
	}

	// 11. Set executable permissions
	if err := KubectlChmod(opts.Namespace, opts.Pod, opts.Container, opts.BinPath, opts.DryRun); err != nil {
		return err
	}

	// 12. Execute netdbg in the pod
	if err := KubectlExec(opts.Namespace, opts.Pod, opts.Container, opts.BinPath, args, opts.DryRun); err != nil {
		return err
	}

	// 13. Cleanup if requested
	if opts.Cleanup {
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
