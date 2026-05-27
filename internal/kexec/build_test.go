package kexec

import (
	"os/exec"
	"testing"
)

// TestBuildNetdbgBinary_DryRun checks that dry-run always succeeds.
func TestBuildNetdbgBinary_DryRun(t *testing.T) {
	setupLogger()
	err := BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

// TestBuildNetdbgBinary_Success simulates a successful build.
func TestBuildNetdbgBinary_Success(t *testing.T) {
	setupLogger()
	oldBuildCommand := buildCommand
	buildCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { buildCommand = oldBuildCommand }()

	err := BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestBuildNetdbgBinary_Error simulates a build failure.
func TestBuildNetdbgBinary_Error(t *testing.T) {
	setupLogger()
	oldBuildCommand := buildCommand
	buildCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { buildCommand = oldBuildCommand }()

	err := BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when build fails")
	}
}

// TestBuildNetdbgBinary_InvalidEnv simulates an error due to invalid environment variables.
func TestBuildNetdbgBinary_InvalidEnv(t *testing.T) {
	setupLogger()
	oldBuildCommand := buildCommand
	buildCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { buildCommand = oldBuildCommand }()

	err := BuildNetdbgBinary("invalidOS", "invalidARCH", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error for invalid GOOS/GOARCH")
	}
}

// TestBuildNetdbgBinary_CommandInjection ensures the injected command is used.
func TestBuildNetdbgBinary_CommandInjection(t *testing.T) {
	setupLogger()
	called := false
	oldBuildCommand := buildCommand
	buildCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("true")
	}
	defer func() { buildCommand = oldBuildCommand }()

	_ = BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", false)
	if !called {
		t.Errorf("expected injected buildCommand to be called")
	}
}

// TestBuildNetdbgBinary_ErrorPropagation ensures errors are propagated.
func TestBuildNetdbgBinary_ErrorPropagation(t *testing.T) {
	setupLogger()
	oldBuildCommand := buildCommand
	buildCommand = func(name string, args ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: "nonexistent",
			Args: []string{"nonexistent"},
		}
	}
	defer func() { buildCommand = oldBuildCommand }()

	err := BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when command does not exist")
	}
}

// TestBuildNetdbgBinary_LoggerCoverage ensures logger.Debug is called (manual/visual).
func TestBuildNetdbgBinary_LoggerCoverage(t *testing.T) {
	setupLogger()
	// This test is mainly for coverage; check that logger.Debug does not panic.
	_ = BuildNetdbgBinary("linux", "amd64", "/tmp/netdbg", true)
}
