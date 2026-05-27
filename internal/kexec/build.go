package kexec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/feliux/netdbg/internal/logger"
)

// buildCommand is a variable for exec.Command, allowing injection in tests.
var buildCommand = exec.Command

// BuildNetdbgBinary compiles netdbg for the specified OS/ARCH and outputs to localBin.
func BuildNetdbgBinary(goos, goarch, localBin string, dryRun bool) error {
	if dryRun {
		logger.Debug("[dry-run] would build netdbg binary", "goos", goos, "goarch", goarch, "output", localBin)
		return nil
	}
	logger.Debug("building netdbg binary", "goos", goos, "goarch", goarch, "output", localBin)
	cmd := buildCommand("go", "build", "-o", localBin, ".")
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Debug("failed to build netdbg binary", "error", err)
		return fmt.Errorf("failed to build netdbg binary: %w", err)
	}
	logger.Debug("successfully built netdbg binary", "output", localBin)
	return nil
}
