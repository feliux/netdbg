package kexec

import (
	"errors"
	"os/exec"
	"testing"
)

// --- CheckKubectlAvailable ---

func TestCheckKubectlAvailable_Success(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "/usr/bin/kubectl", nil }
	defer func() { lookPath = oldLookPath }()
	if err := CheckKubectlAvailable(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

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

func TestCheckGoAvailable_Success(t *testing.T) {
	setupLogger()
	oldLookPath := lookPath
	lookPath = func(file string) (string, error) { return "/usr/bin/go", nil }
	defer func() { lookPath = oldLookPath }()
	if err := CheckGoAvailable(); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

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
