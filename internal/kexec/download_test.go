package kexec

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// mockReadCloser implements io.ReadCloser for testing
type mockReadCloser struct {
	io.Reader
}

func (m *mockReadCloser) Close() error { return nil }

func TestDownloadNetdbgBinary_DryRun(t *testing.T) {
	setupLogger()
	oldHttpGet := httpGet
	defer func() { httpGet = oldHttpGet }()
	err := DownloadNetdbgBinary("https://example.com/netdbg", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

func TestDownloadNetdbgBinary_Success(t *testing.T) {
	setupLogger()
	oldHttpGet := httpGet
	defer func() { httpGet = oldHttpGet }()

	// Mock httpGet to return a simple reader
	httpGet = func(url string) (*http.Response, error) {
		body := &mockReadCloser{Reader: strings.NewReader("netdbg-binary-content")}
		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}

	tmpfile, err := os.CreateTemp("", "netdbg_test_download")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	err = DownloadNetdbgBinary("https://example.com/netdbg", tmpfile.Name(), false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
	// Optionally, check file content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("failed to read downloaded file: %v", err)
	}
	if string(content) != "netdbg-binary-content" {
		t.Errorf("unexpected file content: %s", string(content))
	}
}

func TestDownloadNetdbgBinary_HttpError(t *testing.T) {
	setupLogger()
	oldHttpGet := httpGet
	defer func() { httpGet = oldHttpGet }()
	httpGet = func(url string) (*http.Response, error) {
		return nil, errors.New("network error")
	}
	err := DownloadNetdbgBinary("https://example.com/netdbg", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when httpGet fails")
	}
}

func TestDownloadNetdbgBinary_CreateFileError(t *testing.T) {
	setupLogger()
	oldHttpGet := httpGet
	defer func() { httpGet = oldHttpGet }()
	httpGet = func(url string) (*http.Response, error) {
		body := &mockReadCloser{Reader: strings.NewReader("irrelevant")}
		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}
	// Try to write to an invalid path
	err := DownloadNetdbgBinary("https://example.com/netdbg", "/invalid/path/netdbg", false)
	if err == nil {
		t.Errorf("expected error when file creation fails")
	}
}

func TestDownloadNetdbgBinary_CopyError(t *testing.T) {
	setupLogger()
	oldHttpGet := httpGet
	defer func() { httpGet = oldHttpGet }()
	// Reader that always errors
	errorReader := &mockReadCloser{Reader: errorReader{}}
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       errorReader,
		}, nil
	}
	tmpfile, err := os.CreateTemp("", "netdbg_test_download")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	err = DownloadNetdbgBinary("https://example.com/netdbg", tmpfile.Name(), false)
	if err == nil {
		t.Errorf("expected error when io.Copy fails")
	}
}

// errorReader implements io.Reader and always returns an error
type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}
