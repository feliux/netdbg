package kexec

import (
	"github.com/spf13/pflag"
)

// NetdbgReleaseBaseURL is the default base URL for downloading netdbg releases.
const NetdbgReleaseBaseURL = "https://github.com/feliux/netdbg/releases/download"

// Options holds all configuration for the kexec operation.
type Options struct {
	Namespace      string // Kubernetes namespace
	Pod            string // Pod name
	Container      string // Container name (optional)
	BinPath        string // Path in pod for netdbg binary
	Mode           string // "download" or "build"
	Version        string // netdbg version to download
	Goos           string // Target OS for netdbg binary
	Goarch         string // Target architecture for netdbg binary
	Cleanup        bool   // Remove binary from pod after execution
	KeepLocal      bool   // Keep local binary after operation
	DryRun         bool   // Show actions without executing them
	ReleaseBaseURL string // Optional override for the release base URL
}

// ParseOptionsFromFlags extracts Options from Cobra/pflag.FlagSet.
func ParseOptionsFromFlags(flags *pflag.FlagSet) *Options {
	opts := &Options{}
	opts.Namespace, _ = flags.GetString("namespace")
	opts.Pod, _ = flags.GetString("pod")
	opts.Container, _ = flags.GetString("container")
	opts.BinPath, _ = flags.GetString("bin-path")
	opts.Mode, _ = flags.GetString("mode")
	opts.Version, _ = flags.GetString("version")
	opts.Goos, _ = flags.GetString("goos")
	opts.Goarch, _ = flags.GetString("goarch")
	opts.Cleanup, _ = flags.GetBool("cleanup")
	opts.KeepLocal, _ = flags.GetBool("keep-local")
	opts.DryRun, _ = flags.GetBool("dry-run")
	opts.ReleaseBaseURL, _ = flags.GetString("release-base-url")
	if opts.ReleaseBaseURL == "" {
		opts.ReleaseBaseURL = NetdbgReleaseBaseURL
	}
	return opts
}
