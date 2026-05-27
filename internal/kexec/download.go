package kexec

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/feliux/netdbg/internal/logger"
)

// httpGet is injectable for testing.
var httpGet = http.Get

// DownloadNetdbgBinary downloads the netdbg binary from the given URL to localBin.
func DownloadNetdbgBinary(url, localBin string, dryRun bool) error {
	if dryRun {
		logger.Debug("[dry-run] would download netdbg binary", "url", url, "dest", localBin)
		return nil
	}
	logger.Debug("downloading netdbg binary", "url", url, "dest", localBin)
	resp, err := httpGet(url)
	if err != nil {
		logger.Debug("failed to download netdbg binary", "url", url, "error", err)
		return fmt.Errorf("failed to download netdbg: %v", err)
	}
	defer resp.Body.Close()
	out, err := os.Create(localBin)
	if err != nil {
		logger.Debug("failed to create output file for netdbg binary", "file", localBin, "error", err)
		return fmt.Errorf("failed to create file %s: %v", localBin, err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Debug("failed to write netdbg binary to file", "file", localBin, "error", err)
		return fmt.Errorf("failed to write netdbg binary: %v", err)
	}
	logger.Debug("netdbg binary downloaded successfully", "file", localBin)
	return nil
}
