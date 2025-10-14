package kexec

import (
	"os/exec"
	"testing"
)

func TestKubectlCp_DryRun(t *testing.T) {
	err := KubectlCp("default", "mypod", "", "/tmp/netdbg", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

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

func TestKubectlChmod_DryRun(t *testing.T) {
	err := KubectlChmod("default", "mypod", "", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

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

func TestKubectlExec_DryRun(t *testing.T) {
	err := KubectlExec("default", "mypod", "", "/tmp/netdbg", []string{"nc", "-a", "8.8.8.8"}, true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

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

func TestKubectlCleanup_DryRun(t *testing.T) {
	err := KubectlCleanup("default", "mypod", "", "/tmp/netdbg", true)
	if err != nil {
		t.Errorf("expected no error in dry-run, got: %v", err)
	}
}

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
