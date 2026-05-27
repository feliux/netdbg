package kexec

import (
	"errors"
	"os/exec"
	"testing"
)

// --- CheckKubectlAvailable ---

// TestCheckKubectlAvailable_Success verifies success when kubectl is available.
func TestCheckKubectlAvailable_Success(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "/usr/bin/kubectl", nil }
	defer func() { lookPath = oldLookPath }()
	if err := CheckKubectlAvailable(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckKubectlAvailable_Error verifies failure when kubectl is missing.
func TestCheckKubectlAvailable_Error(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "", errors.New("not found") }
	defer func() { lookPath = oldLookPath }()
	if err := CheckKubectlAvailable(); err == nil {
		t.Errorf("expected error when kubectl is not found")
	}
}

// --- CheckKubeCluster ---

// TestCheckKubeCluster_Success verifies cluster connectivity succeeds.
func TestCheckKubeCluster_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { command = oldCommand }()
	if err := CheckKubeCluster(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckKubeCluster_Error verifies cluster connectivity failures.
func TestCheckKubeCluster_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	if err := CheckKubeCluster(); err == nil {
		t.Errorf("expected error when cluster-info fails")
	}
}

// --- CheckPodRunning ---

// TestCheckPodRunning_Success verifies running pods are accepted.
func TestCheckPodRunning_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "Running")
	}
	defer func() { command = oldCommand }()
	if err := CheckPodRunning("default", "mypod"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckPodRunning_NotRunning verifies non-running pods are rejected.
func TestCheckPodRunning_NotRunning(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "Pending")
	}
	defer func() { command = oldCommand }()
	if err := CheckPodRunning("default", "mypod"); err == nil {
		t.Errorf("expected error when pod is not Running")
	}
}

// TestCheckPodRunning_Error verifies errors when pod status cannot be read.
func TestCheckPodRunning_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	if err := CheckPodRunning("default", "mypod"); err == nil {
		t.Errorf("expected error when kubectl get pod fails")
	}
}

// --- CheckExecPermission ---

// TestCheckExecPermission_Success verifies exec permission checks pass.
func TestCheckExecPermission_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { command = oldCommand }()
	if err := CheckExecPermission("default", "mypod", ""); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckExecPermission_Error verifies exec permission checks fail.
func TestCheckExecPermission_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	if err := CheckExecPermission("default", "mypod", ""); err == nil {
		t.Errorf("expected error when exec is not permitted")
	}
}

// --- CheckWritable ---

// TestCheckWritable_Success verifies write permission checks pass.
func TestCheckWritable_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { command = oldCommand }()
	if err := CheckWritable("default", "mypod", "", "/tmp"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckWritable_Error verifies write permission checks fail.
func TestCheckWritable_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	if err := CheckWritable("default", "mypod", "", "/tmp"); err == nil {
		t.Errorf("expected error when path is not writable")
	}
}

// --- CheckBinNotExists ---

// TestCheckBinNotExists_Success verifies absence of a binary is accepted.
func TestCheckBinNotExists_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { command = oldCommand }()
	if err := CheckBinNotExists("default", "mypod", "", "/tmp/netdbg"); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckBinNotExists_Error verifies existing binaries are detected.
func TestCheckBinNotExists_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	if err := CheckBinNotExists("default", "mypod", "", "/tmp/netdbg"); err == nil {
		t.Errorf("expected error when binary exists")
	}
}

// --- CheckGoAvailable ---

// TestCheckGoAvailable_Success verifies Go availability checks pass.
func TestCheckGoAvailable_Success(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "/usr/bin/go", nil }
	defer func() { lookPath = oldLookPath }()
	if err := CheckGoAvailable(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

// TestCheckGoAvailable_Error verifies Go availability checks fail.
func TestCheckGoAvailable_Error(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "", errors.New("not found") }
	defer func() { lookPath = oldLookPath }()
	if err := CheckGoAvailable(); err == nil {
		t.Errorf("expected error when go is not found")
	}
}

// --- GetPodArchOS ---

// TestGetPodArchOS_Success verifies pod architecture/OS detection succeeds.
func TestGetPodArchOS_Success(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		if len(arg) > 0 && arg[len(arg)-2] == "uname" && arg[len(arg)-1] == "-m" {
			return exec.Command("echo", "x86_64")
		}
		if len(arg) > 0 && arg[len(arg)-2] == "uname" && arg[len(arg)-1] == "-s" {
			return exec.Command("echo", "Linux")
		}
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	arch, osname, err := GetPodArchOS("default", "mypod", "")
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
	if arch != "x86_64" || osname != "Linux" {
		t.Errorf("unexpected arch/os: %s/%s", arch, osname)
	}
}

// TestGetPodArchOS_Error verifies pod architecture/OS detection failures.
func TestGetPodArchOS_Error(t *testing.T) {
	setupLogger()
	oldCommand := command
	command = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() { command = oldCommand }()
	_, _, err := GetPodArchOS("default", "mypod", "")
	if err == nil {
		t.Errorf("expected error when uname fails")
	}
}
