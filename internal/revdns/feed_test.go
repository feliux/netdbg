package revdns

import (
	"fmt"
	"os"
	"testing"
)

// TestFeedFromFile checks that feedFromFile sends all lines to the channel.
func TestFeedFromFile(t *testing.T) {
	lines := []string{"1.1.1.1", "8.8.8.8", "127.0.0.1"}
	tmpfile, err := os.CreateTemp("", "feedfromfile_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	for _, line := range lines {
		fmt.Fprintln(tmpfile, line)
	}
	tmpfile.Close()

	ch := make(chan string, len(lines))
	err = feedFromFile(tmpfile.Name(), ch)
	if err != nil {
		t.Fatalf("feedFromFile returned error: %v", err)
	}
	close(ch)

	var got []string
	for ip := range ch {
		got = append(got, ip)
	}
	// Since the channel is buffered, we can range after close
	if len(got) != len(lines) {
		t.Errorf("expected %d lines, got %d", len(lines), len(got))
	}
	for _, want := range lines {
		found := false
		for _, have := range got {
			if have == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find %q in output", want)
		}
	}
}
