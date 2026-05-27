package kexec

import (
	"os/exec"
	"testing"
)

// TestKubectlCp_DryRun verifies dry-run behavior for kubectl cp.
func TestKubectlCp_DryRun(t *testing.T) {
	err := KubectlCp("default", "mypod", "", "/tmp/netdbg", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

// TestKubectlCp_Success verifies a successful kubectl cp execution.
func TestKubectlCp_Success(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlCp("default", "mypod", "", "/tmp/netdbg", "/tmp/netdbg", false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestKubectlCp_Error verifies kubectl cp errors are reported.
func TestKubectlCp_Error(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlCp("default", "mypod", "", "/tmp/netdbg", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when kubectl cp fails")
	}
}

// --- KubectlChmod ---

// TestKubectlChmod_DryRun verifies dry-run behavior for kubectl chmod.
func TestKubectlChmod_DryRun(t *testing.T) {
	err := KubectlChmod("default", "mypod", "", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

// TestKubectlChmod_Success verifies a successful kubectl chmod execution.
func TestKubectlChmod_Success(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlChmod("default", "mypod", "", "/tmp/netdbg", false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestKubectlChmod_Error verifies kubectl chmod errors are reported.
func TestKubectlChmod_Error(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlChmod("default", "mypod", "", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when kubectl chmod fails")
	}
}

// --- KubectlExec ---

// TestKubectlExec_DryRun verifies dry-run behavior for kubectl exec.
func TestKubectlExec_DryRun(t *testing.T) {
	err := KubectlExec("default", "mypod", "", "/tmp/netdbg", []string{"nc", "-a", "8.8.8.8"}, true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

// TestKubectlExec_Success verifies a successful kubectl exec execution.
func TestKubectlExec_Success(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlExec("default", "mypod", "", "/tmp/netdbg", []string{"nc", "-a", "8.8.8.8"}, false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestKubectlExec_Error verifies kubectl exec errors are reported.
func TestKubectlExec_Error(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlExec("default", "mypod", "", "/tmp/netdbg", []string{"nc", "-a", "8.8.8.8"}, false)
	if err == nil {
		t.Errorf("expected error when kubectl exec fails")
	}
}

// --- KubectlCleanup ---

// TestKubectlCleanup_DryRun verifies dry-run behavior for kubectl cleanup.
func TestKubectlCleanup_DryRun(t *testing.T) {
	err := KubectlCleanup("default", "mypod", "", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

// TestKubectlCleanup_Success verifies a successful kubectl cleanup execution.
func TestKubectlCleanup_Success(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlCleanup("default", "mypod", "", "/tmp/netdbg", false)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestKubectlCleanup_Error verifies kubectl cleanup errors are reported.
func TestKubectlCleanup_Error(t *testing.T) {
	oldCmd := kubectlCommand
	kubectlCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { kubectlCommand = oldCmd }()
	err := KubectlCleanup("default", "mypod", "", "/tmp/netdbg", false)
	if err == nil {
		t.Errorf("expected error when kubectl cleanup fails")
	}
}
