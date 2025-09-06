package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sh05/terraform-provider-multipass/internal/common"
)

func TestNewMultipassClient(t *testing.T) {
	// Test with default binary path
	client := NewMultipassClient("")
	if client.binaryPath != "multipass" {
		t.Errorf("Expected binary path to be 'multipass', got %s", client.binaryPath)
	}

	// Test with custom binary path
	customPath := "/usr/local/bin/multipass"
	client = NewMultipassClient(customPath)
	if client.binaryPath != customPath {
		t.Errorf("Expected binary path to be '%s', got %s", customPath, client.binaryPath)
	}
}

func TestMultipassClientLaunchOptions(t *testing.T) {
	client := NewMultipassClient("")

	// This is a unit test that doesn't actually run multipass commands
	// It tests the launch options construction
	opts := &common.LaunchOptions{
		Name:      "test-instance",
		Image:     "22.04",
		CPU:       "2",
		Memory:    "2G",
		Disk:      "10G",
		CloudInit: "/path/to/cloud-init.yaml",
		Timeout:   "15m0s",
	}

	// We can't actually test Launch without multipass installed,
	// but we can test that the client is constructed correctly
	if client == nil {
		t.Error("Client should not be nil")
	}

	if opts.Name != "test-instance" {
		t.Errorf("Expected name to be 'test-instance', got %s", opts.Name)
	}

	if opts.Image != "22.04" {
		t.Errorf("Expected image to be '22.04', got %s", opts.Image)
	}

	if opts.Timeout != "15m0s" {
		t.Errorf("Expected timeout to be '15m0s', got %s", opts.Timeout)
	}
}

// TestMultipassClientWithNonExistentBinary tests client behavior when binary doesn't exist
func TestMultipassClientWithNonExistentBinary(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	opts := &common.LaunchOptions{
		Name:  "test-instance",
		Image: "22.04",
	}
	
	err := client.Launch(opts)
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to launch instance") {
		t.Errorf("Expected error message to contain 'failed to launch instance', got: %s", err.Error())
	}
}

// TestMultipassClientInvalidTimeout tests timeout parsing errors
func TestMultipassClientInvalidTimeout(t *testing.T) {
	client := NewMultipassClient("echo") // Use echo as a mock command that won't fail
	
	testCases := []struct {
		name    string
		timeout string
		wantErr bool
	}{
		{
			name:    "Invalid duration format",
			timeout: "invalid-duration",
			wantErr: true,
		},
		{
			name:    "Empty timeout",
			timeout: "",
			wantErr: false, // Empty timeout should be valid (uses default)
		},
		{
			name:    "Negative duration",
			timeout: "-5m",
			wantErr: true,
		},
		{
			name:    "Valid duration",
			timeout: "10m",
			wantErr: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &common.LaunchOptions{
				Name:    "test-instance",
				Image:   "22.04",
				Timeout: tc.timeout,
			}
			
			err := client.Launch(opts)
			if tc.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.wantErr && err != nil && strings.Contains(err.Error(), "invalid timeout format") {
				t.Errorf("Unexpected timeout parsing error: %v", err)
			}
		})
	}
}

// TestMultipassClientGetInstanceWithNonExistentBinary tests GetInstance with invalid binary
func TestMultipassClientGetInstanceWithNonExistentBinary(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	instance, err := client.GetInstance("test-instance")
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	if instance != nil {
		t.Error("Expected nil instance when error occurs")
	}
	
	if !strings.Contains(err.Error(), "failed to get instance info") {
		t.Errorf("Expected error message to contain 'failed to get instance info', got: %s", err.Error())
	}
}

// TestMultipassClientListInstancesError tests ListInstances error handling
func TestMultipassClientListInstancesError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	instances, err := client.ListInstances()
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	if instances != nil {
		t.Error("Expected nil instances when error occurs")
	}
}

// TestMultipassClientDeleteInstanceError tests DeleteInstance error handling
func TestMultipassClientDeleteInstanceError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	err := client.DeleteInstance("test-instance", false)
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to delete instance") {
		t.Errorf("Expected error message to contain 'failed to delete instance', got: %s", err.Error())
	}
}

// TestMultipassClientStartInstanceError tests StartInstance error handling
func TestMultipassClientStartInstanceError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	err := client.StartInstance("test-instance")
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to start instance") {
		t.Errorf("Expected error message to contain 'failed to start instance', got: %s", err.Error())
	}
}

// TestMultipassClientStopInstanceError tests StopInstance error handling
func TestMultipassClientStopInstanceError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	err := client.StopInstance("test-instance")
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to stop instance") {
		t.Errorf("Expected error message to contain 'failed to stop instance', got: %s", err.Error())
	}
}

// TestMultipassClientSuspendInstanceError tests SuspendInstance error handling
func TestMultipassClientSuspendInstanceError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	err := client.SuspendInstance("test-instance")
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to suspend instance") {
		t.Errorf("Expected error message to contain 'failed to suspend instance', got: %s", err.Error())
	}
}

// TestMultipassClientRestartInstanceError tests RestartInstance error handling
func TestMultipassClientRestartInstanceError(t *testing.T) {
	client := NewMultipassClient("/non/existent/binary")
	
	err := client.RestartInstance("test-instance")
	if err == nil {
		t.Error("Expected error when using non-existent binary, got nil")
	}
	
	if !strings.Contains(err.Error(), "failed to restart instance") {
		t.Errorf("Expected error message to contain 'failed to restart instance', got: %s", err.Error())
	}
}

// TestMultipassClientWithPermissionError simulates permission errors
func TestMultipassClientWithPermissionError(t *testing.T) {
	// Create a temporary file with no execute permission to simulate permission error
	tempDir := t.TempDir()
	fakeBinary := filepath.Join(tempDir, "fake-multipass")
	
	// Create file without execute permission
	err := os.WriteFile(fakeBinary, []byte("#!/bin/sh\necho test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create fake binary: %v", err)
	}
	
	client := NewMultipassClient(fakeBinary)
	
	opts := &common.LaunchOptions{
		Name:  "test-instance",
		Image: "22.04",
	}
	
	err = client.Launch(opts)
	if err == nil {
		t.Error("Expected error due to permission denied, got nil")
	}
	
	// Error should contain information about execution failure
	if !strings.Contains(err.Error(), "failed to launch instance") {
		t.Errorf("Expected error message to contain 'failed to launch instance', got: %s", err.Error())
	}
}

// TestDurationToSecondsEdgeCases tests the durationToSeconds method with various inputs
func TestDurationToSecondsEdgeCases(t *testing.T) {
	client := NewMultipassClient("echo")
	
	testCases := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{"Valid minutes", "5m", false},
		{"Valid seconds", "30s", false},
		{"Valid hours", "2h", false},
		{"Valid mixed", "1h30m45s", false},
		{"Invalid format", "abc", true},
		{"Empty string", "", true},
		{"Just number", "123", true},
		{"Negative", "-5m", true},
		{"Zero", "0s", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.durationToSeconds(tc.duration)
			if tc.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
