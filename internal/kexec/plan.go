package kexec

import (
	"fmt"
	"strings"
)

// DryRunPlan returns a human-readable list of dry-run steps.
func DryRunPlan(opts *Options, args []string) []string {
	localBin := "./netdbg"
	target := formatTarget(opts)
	steps := []string{
		"dry-run: planned actions (nothing will be executed)",
		fmt.Sprintf("- target: %s", target),
		fmt.Sprintf("- mode: %s", opts.Mode),
		fmt.Sprintf("- local binary: %s", localBin),
		fmt.Sprintf("- pod binary path: %s", opts.BinPath),
		fmt.Sprintf("- target platform: %s/%s", opts.Goos, opts.Goarch),
	}

	if opts.Mode == "download" {
		steps = append(steps, fmt.Sprintf("- version: %s", opts.Version))
		if opts.ReleaseBaseURL != "" {
			steps = append(steps, fmt.Sprintf("- release base URL: %s", opts.ReleaseBaseURL))
		}
	}

	steps = append(steps,
		"- check kubectl availability",
		"- check cluster connectivity",
		fmt.Sprintf("- check pod is running (%s)", target),
		fmt.Sprintf("- check exec permission (%s)", target),
		fmt.Sprintf("- check write permission on %s (%s)", opts.BinPath, target),
		fmt.Sprintf("- check no binary exists at %s (%s)", opts.BinPath, target),
	)

	if opts.Mode == "build" {
		steps = append(steps, fmt.Sprintf("- build netdbg (GOOS=%s GOARCH=%s) -> %s", opts.Goos, opts.Goarch, localBin))
	} else {
		binSource := fmt.Sprintf("%s/%s/netdbg-%s-%s", opts.ReleaseBaseURL, opts.Version, opts.Goos, opts.Goarch)
		steps = append(steps, fmt.Sprintf("- download netdbg (%s) -> %s", binSource, localBin))
	}

	steps = append(steps,
		fmt.Sprintf("- copy binary to pod: %s (%s)", opts.BinPath, target),
		fmt.Sprintf("- set executable permission on %s (%s)", opts.BinPath, target),
	)

	execArgs := strings.Join(args, " ")
	if execArgs == "" {
		execArgs = "(no args)"
	}
	steps = append(steps, fmt.Sprintf("- execute in pod: %s %s (%s)", opts.BinPath, execArgs, target))

	if opts.Cleanup {
		steps = append(steps, fmt.Sprintf("- remove binary from pod: %s (%s)", opts.BinPath, target))
	} else {
		steps = append(steps, fmt.Sprintf("- keep binary in pod: %s (%s)", opts.BinPath, target))
	}
	if !opts.KeepLocal {
		steps = append(steps, fmt.Sprintf("- remove local binary: %s", localBin))
	} else {
		steps = append(steps, fmt.Sprintf("- keep local binary: %s", localBin))
	}

	return steps
}

// formatTarget builds a short description of the Kubernetes target.
func formatTarget(opts *Options) string {
	if opts.Container == "" {
		return fmt.Sprintf("namespace=%s, pod=%s", opts.Namespace, opts.Pod)
	}
	return fmt.Sprintf("namespace=%s, pod=%s, container=%s", opts.Namespace, opts.Pod, opts.Container)
}
