package cmd

import (
	"github.com/feliux/netdbg/internal/kexec"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(kexecCmd)
	kexecCmd.Flags().StringP("namespace", "n", "default", "pod namespace")
	kexecCmd.Flags().StringP("pod", "p", "", "pod name (required)")
	kexecCmd.Flags().StringP("container", "c", "", "container name (optional)")
	kexecCmd.Flags().String("bin-path", "/tmp/netdbg", "path in pod for netdbg binary")
	kexecCmd.Flags().String("mode", "download", "how to obtain netdbg binary: 'download' (from GitHub release) or 'build' (local compile)")
	kexecCmd.Flags().String("version", "latest", "netdbg version to download (default: latest)")
	kexecCmd.Flags().String("goos", "linux", "target os for netdbg binary")
	kexecCmd.Flags().String("goarch", "amd64", "target architecture for netdbg binary")
	kexecCmd.Flags().Bool("cleanup", false, "remove binary from pod after execution")
	kexecCmd.Flags().Bool("keep-local", false, "keep local binary after operation")
	kexecCmd.Flags().Bool("dry-run", false, "show actions without executing them")
	kexecCmd.Flags().String("release-base-url", "", "override the default netdbg release base URL")
	kexecCmd.MarkFlagRequired("pod")
}

var kexecCmd = &cobra.Command{
	Use:   "kexec",
	Short: "Run netdbg commands inside a Kubernetes pod (auto-upload and exec)",
	Long: `Run netdbg commands inside a Kubernetes pod, automatically handling binary upload and execution.

Usage examples:
  netdbg kexec -n myns -p mypod --mode download --version v1.2.3 -- nc -a 8.8.8.8 -p 53
  netdbg kexec -n myns -p mypod --mode build --goarch arm64 -- revdns -a 8.8.8.8
  netdbg kexec -n myns -p mypod --cleanup -- nc -a 1.1.1.1 -p 80
`,
	DisableFlagParsing: false,
	Run: func(cmd *cobra.Command, args []string) {
		kexec.ExecuteCommand(cmd, args)
	},
}
