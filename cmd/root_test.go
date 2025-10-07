package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmdConfig(t *testing.T) {
	// Verify that the root command is configured correctly
	if rootCmd.Use != "netdbg" {
		t.Errorf("Expected rootCmd.Use to be 'netdbg', got '%s'", rootCmd.Use)
	}
	if rootCmd.Short != "Net debugger CLI" {
		t.Errorf("Expected rootCmd.Short to be 'Net debugger CLI', got '%s'", rootCmd.Short)
	}
	if rootCmd.Long != "Set of tools for testing and debugging connectivity issues." {
		t.Errorf("Expected rootCmd.Long to be 'Set of tools for testing and debugging connectivity issues.', got '%s'", rootCmd.Long)
	}
}

func TestRootCmdExecution(t *testing.T) {
	// Capture the output of the root command
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Execute the root command without arguments
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Expected rootCmd.Execute() to run without error, got: %v", err)
	}

	// Verify that there is no unexpected output
	// output := buf.String()
	// if output != "" {
	// 	t.Errorf("Expected no output from rootCmd.Execute(), got: %s", output)
	// }
}

func TestRootCmdWithInvalidArgs(t *testing.T) {
	// Capture the output of the root command
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Execute the root command with invalid arguments
	rootCmd.SetArgs([]string{"invalid"})
	err := rootCmd.Execute()
	if err == nil {
		t.Errorf("Expected rootCmd.Execute() to return an error for invalid arguments")
	}
}
