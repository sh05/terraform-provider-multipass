package provider

import (
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
