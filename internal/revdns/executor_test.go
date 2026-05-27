package revdns

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TestDefaultExecutor_Execute_Success verifies successful reverse DNS execution.
func TestDefaultExecutor_Execute_Success(t *testing.T) {
	setupLogger()
	opts := &Options{
		Addr:       "8.8.8.8", // Google DNS, should resolve to dns.google
		Threads:    1,
		DomainOnly: false,
	}
	executor := &DefaultExecutor{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var gotDomain string
	var gotIP string
	for result := range executor.Execute(ctx, opts) {
		if result.Error != nil {
			t.Fatalf("unexpected error: %v", result.Error)
		}
		gotDomain = result.Domain
		gotIP = result.IP
	}

	if gotIP != "8.8.8.8" {
		t.Errorf("expected IP 8.8.8.8, got %s", gotIP)
	}
	if gotDomain == "" {
		t.Errorf("expected a domain for 8.8.8.8, got empty string")
	}
}

// TestDefaultExecutor_Execute_Fail verifies failures for unresolvable IPs.
func TestDefaultExecutor_Execute_Fail(t *testing.T) {
	setupLogger()
	opts := &Options{
		Addr:       "192.0.2.123", // TEST-NET-1, should not resolve
		Threads:    1,
		DomainOnly: false,
	}
	executor := &DefaultExecutor{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var gotError error
	for result := range executor.Execute(ctx, opts) {
		if result.Error != nil {
			gotError = result.Error
		}
	}
	if gotError == nil {
		t.Errorf("expected error for unresolvable IP, got nil")
	}
}

// TestDefaultExecutor_Execute_DomainOnly verifies domain-only output behavior.
func TestDefaultExecutor_Execute_DomainOnly(t *testing.T) {
	setupLogger()
	opts := &Options{
		Addr:       "8.8.8.8",
		Threads:    1,
		DomainOnly: true,
	}
	executor := &DefaultExecutor{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var gotDomain string
	for result := range executor.Execute(ctx, opts) {
		if result.Error != nil {
			t.Fatalf("unexpected error: %v", result.Error)
		}
		gotDomain = result.Domain
	}
	if gotDomain == "" {
		t.Errorf("expected a domain for 8.8.8.8, got empty string")
	}
}

// TestDefaultExecutor_Execute_FromFile verifies file-driven reverse DNS execution.
func TestDefaultExecutor_Execute_FromFile(t *testing.T) {
	setupLogger()
	// Prepare a temporary file with two IPs: one resolvable, one not
	tmpfile, err := os.CreateTemp("", "revdns_test_ips")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	ips := []string{"8.8.8.8", "192.0.2.123"}
	for _, ip := range ips {
		fmt.Fprintln(tmpfile, ip)
	}
	tmpfile.Close()

	opts := &Options{
		File:       tmpfile.Name(),
		Threads:    2,
		DomainOnly: false,
	}
	executor := &DefaultExecutor{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	foundGood := false
	foundBad := false
	for result := range executor.Execute(ctx, opts) {
		if result.Error != nil {
			if strings.Contains(result.IP, "192.0.2.123") {
				foundBad = true
			}
			continue
		}
		if result.IP == "8.8.8.8" && result.Domain != "" {
			foundGood = true
		}
	}
	if !foundGood {
		t.Errorf("expected to resolve 8.8.8.8 from file")
	}
	if !foundBad {
		t.Errorf("expected error for 192.0.2.123 from file")
	}
}
