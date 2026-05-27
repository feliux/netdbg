package revdns

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

// feedFromFile reads IPs from a file and sends them to the work channel, one per line.
func feedFromFile(filePath string, work chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("failed to close file", "err", cerr)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		work <- scanner.Text()
	}
	return scanner.Err()
}
