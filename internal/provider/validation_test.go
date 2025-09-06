package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/sh05/terraform-provider-multipass/internal/common"
)

// TestValidateInstanceName tests instance name validation
func TestValidateInstanceName(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid simple name",
			input:   "test-instance",
			wantErr: false,
		},
		{
			name:    "Valid name with numbers",
			input:   "test-instance-123",
			wantErr: false,
		},
		{
			name:    "Valid name with underscores",
			input:   "test_instance_name",
			wantErr: false,
		},
		{
			name:    "Empty name",
			input:   "",
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name:    "Name with special characters",
			input:   "test@instance#name",
			wantErr: true,
			errMsg:  "invalid characters",
		},
		{
			name:    "Name with spaces",
			input:   "test instance name",
			wantErr: true,
			errMsg:  "spaces not allowed",
		},
		{
			name:    "Name starting with number",
			input:   "123-test-instance",
			wantErr: true,
			errMsg:  "cannot start with number",
		},
		{
			name:    "Name starting with dash",
			input:   "-test-instance",
			wantErr: true,
			errMsg:  "cannot start with dash",
		},
		{
			name:    "Name ending with dash",
			input:   "test-instance-",
			wantErr: true,
			errMsg:  "cannot end with dash",
		},
		{
			name:    "Name too long",
			input:   strings.Repeat("a", 256),
			wantErr: true,
			errMsg:  "name too long",
		},
		{
			name:    "Name with unicode characters",
			input:   "test-❤️-instance",
			wantErr: true,
			errMsg:  "invalid characters",
		},
		{
			name:    "Single character name",
			input:   "a",
			wantErr: false,
		},
		{
			name:    "Name with consecutive dashes",
			input:   "test--instance",
			wantErr: true,
			errMsg:  "consecutive dashes not allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateInstanceName(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateCPUValue tests CPU value validation
func TestValidateCPUValue(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid single CPU",
			input:   "1",
			wantErr: false,
		},
		{
			name:    "Valid multiple CPUs",
			input:   "4",
			wantErr: false,
		},
		{
			name:    "Valid maximum CPUs",
			input:   "16",
			wantErr: false,
		},
		{
			name:    "Zero CPUs",
			input:   "0",
			wantErr: true,
			errMsg:  "must be greater than 0",
		},
		{
			name:    "Negative CPUs",
			input:   "-1",
			wantErr: true,
			errMsg:  "must be greater than 0",
		},
		{
			name:    "Non-numeric CPU",
			input:   "abc",
			wantErr: true,
			errMsg:  "must be a number",
		},
		{
			name:    "Decimal CPU",
			input:   "2.5",
			wantErr: true,
			errMsg:  "must be an integer",
		},
		{
			name:    "Excessive CPUs",
			input:   "9999",
			wantErr: true,
			errMsg:  "exceeds maximum",
		},
		{
			name:    "Empty CPU value",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCPUValue(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateMemoryValue tests memory value validation
func TestValidateMemoryValue(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid memory in MB",
			input:   "512M",
			wantErr: false,
		},
		{
			name:    "Valid memory in GB",
			input:   "2G",
			wantErr: false,
		},
		{
			name:    "Valid memory in bytes",
			input:   "1073741824",
			wantErr: false,
		},
		{
			name:    "Memory too small",
			input:   "128M",
			wantErr: true,
			errMsg:  "minimum memory",
		},
		{
			name:    "Memory too large",
			input:   "999G",
			wantErr: true,
			errMsg:  "exceeds maximum",
		},
		{
			name:    "Invalid memory format",
			input:   "abc",
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name:    "Negative memory",
			input:   "-1G",
			wantErr: true,
			errMsg:  "must be positive",
		},
		{
			name:    "Zero memory",
			input:   "0G",
			wantErr: true,
			errMsg:  "must be greater than 0",
		},
		{
			name:    "Memory with invalid unit",
			input:   "1X",
			wantErr: true,
			errMsg:  "invalid unit",
		},
		{
			name:    "Empty memory value",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "Decimal memory",
			input:   "1.5G",
			wantErr: true,
			errMsg:  "decimal values not supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMemoryValue(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateDiskValue tests disk value validation
func TestValidateDiskValue(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid disk in GB",
			input:   "10G",
			wantErr: false,
		},
		{
			name:    "Valid disk in MB",
			input:   "1024M",
			wantErr: false,
		},
		{
			name:    "Valid minimum disk",
			input:   "5G",
			wantErr: false,
		},
		{
			name:    "Disk too small",
			input:   "1G",
			wantErr: true,
			errMsg:  "minimum disk",
		},
		{
			name:    "Disk too large",
			input:   "10000G",
			wantErr: true,
			errMsg:  "exceeds maximum",
		},
		{
			name:    "Invalid disk format",
			input:   "invalid",
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name:    "Negative disk",
			input:   "-5G",
			wantErr: true,
			errMsg:  "must be positive",
		},
		{
			name:    "Zero disk",
			input:   "0G",
			wantErr: true,
			errMsg:  "must be greater than 0",
		},
		{
			name:    "Empty disk value",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDiskValue(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateCloudInitFile tests cloud-init file path validation
func TestValidateCloudInitFile(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Empty path (optional)",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Valid YAML file",
			input:   "/path/to/cloud-init.yaml",
			wantErr: false,
		},
		{
			name:    "Valid YML file",
			input:   "/path/to/cloud-init.yml",
			wantErr: false,
		},
		{
			name:    "Relative path",
			input:   "./cloud-init.yaml",
			wantErr: false,
		},
		{
			name:    "Non-YAML file",
			input:   "/path/to/file.txt",
			wantErr: true,
			errMsg:  "must be a YAML file",
		},
		{
			name:    "File with no extension",
			input:   "/path/to/cloud-init",
			wantErr: true,
			errMsg:  "must have .yaml or .yml extension",
		},
		{
			name:    "Invalid characters in path",
			input:   "/path/with\x00null/cloud-init.yaml",
			wantErr: true,
			errMsg:  "invalid characters",
		},
		{
			name:    "Path too long",
			input:   "/" + strings.Repeat("a", 1000) + "/cloud-init.yaml",
			wantErr: true,
			errMsg:  "path too long",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCloudInitFile(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateTimeoutValue tests timeout value validation
func TestValidateTimeoutValue(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid seconds",
			input:   "30s",
			wantErr: false,
		},
		{
			name:    "Valid minutes",
			input:   "5m",
			wantErr: false,
		},
		{
			name:    "Valid hours",
			input:   "1h",
			wantErr: false,
		},
		{
			name:    "Valid mixed duration",
			input:   "1h30m45s",
			wantErr: false,
		},
		{
			name:    "Empty timeout (uses default)",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "abc",
			wantErr: true,
			errMsg:  "invalid duration",
		},
		{
			name:    "Negative timeout",
			input:   "-5m",
			wantErr: true,
			errMsg:  "must be positive",
		},
		{
			name:    "Zero timeout",
			input:   "0s",
			wantErr: true,
			errMsg:  "must be greater than 0",
		},
		{
			name:    "Timeout too large",
			input:   "999h",
			wantErr: true,
			errMsg:  "exceeds maximum",
		},
		{
			name:    "Just a number",
			input:   "123",
			wantErr: true,
			errMsg:  "must include time unit",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTimeoutValue(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateLaunchOptions tests complete launch options validation
func TestValidateLaunchOptions(t *testing.T) {
	testCases := []struct {
		name    string
		opts    *common.LaunchOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid options",
			opts: &common.LaunchOptions{
				Name:   "test-instance",
				Image:  "22.04",
				CPU:    "2",
				Memory: "2G",
				Disk:   "10G",
			},
			wantErr: false,
		},
		{
			name: "Invalid name",
			opts: &common.LaunchOptions{
				Name:   "invalid@name",
				Image:  "22.04",
				CPU:    "2",
				Memory: "2G",
				Disk:   "10G",
			},
			wantErr: true,
			errMsg:  "invalid name",
		},
		{
			name: "Invalid CPU",
			opts: &common.LaunchOptions{
				Name:   "test-instance",
				Image:  "22.04",
				CPU:    "-1",
				Memory: "2G",
				Disk:   "10G",
			},
			wantErr: true,
			errMsg:  "invalid CPU",
		},
		{
			name: "Invalid memory",
			opts: &common.LaunchOptions{
				Name:   "test-instance",
				Image:  "22.04",
				CPU:    "2",
				Memory: "invalid",
				Disk:   "10G",
			},
			wantErr: true,
			errMsg:  "invalid memory",
		},
		{
			name: "Invalid disk",
			opts: &common.LaunchOptions{
				Name:   "test-instance",
				Image:  "22.04",
				CPU:    "2",
				Memory: "2G",
				Disk:   "1G",
			},
			wantErr: true,
			errMsg:  "invalid disk",
		},
		{
			name: "Invalid cloud-init",
			opts: &common.LaunchOptions{
				Name:      "test-instance",
				Image:     "22.04",
				CPU:       "2",
				Memory:    "2G",
				Disk:      "10G",
				CloudInit: "/invalid/file.txt",
			},
			wantErr: true,
			errMsg:  "invalid cloud-init",
		},
		{
			name: "Invalid timeout",
			opts: &common.LaunchOptions{
				Name:    "test-instance",
				Image:   "22.04",
				CPU:     "2",
				Memory:  "2G",
				Disk:    "10G",
				Timeout: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid timeout",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLaunchOptions(tc.opts)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Mock validation functions - these would typically be implemented in the actual provider code

func validateInstanceName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 255 {
		return fmt.Errorf("name too long")
	}
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9\-_]`, name); matched {
		return fmt.Errorf("invalid characters in name")
	}
	if matched, _ := regexp.MatchString(`^[0-9\-]`, name); matched {
		return fmt.Errorf("name cannot start with number or dash")
	}
	if strings.HasSuffix(name, "-") {
		return fmt.Errorf("name cannot end with dash")
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("consecutive dashes not allowed")
	}
	return nil
}

func validateCPUValue(cpu string) error {
	if cpu == "" {
		return fmt.Errorf("CPU value cannot be empty")
	}
	// Add validation logic here
	return nil
}

func validateMemoryValue(memory string) error {
	if memory == "" {
		return fmt.Errorf("memory value cannot be empty")
	}
	// Add validation logic here
	return nil
}

func validateDiskValue(disk string) error {
	if disk == "" {
		return fmt.Errorf("disk value cannot be empty")
	}
	// Add validation logic here
	return nil
}

func validateCloudInitFile(path string) error {
	if path == "" {
		return nil // Optional field
	}
	if len(path) > 1000 {
		return fmt.Errorf("path too long")
	}
	// Add more validation logic here
	return nil
}

func validateTimeoutValue(timeout string) error {
	if timeout == "" {
		return nil // Optional field
	}
	// Add validation logic here
	return nil
}

func validateLaunchOptions(opts *common.LaunchOptions) error {
	if err := validateInstanceName(opts.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if err := validateCPUValue(opts.CPU); err != nil {
		return fmt.Errorf("invalid CPU: %w", err)
	}
	if err := validateMemoryValue(opts.Memory); err != nil {
		return fmt.Errorf("invalid memory: %w", err)
	}
	if err := validateDiskValue(opts.Disk); err != nil {
		return fmt.Errorf("invalid disk: %w", err)
	}
	if err := validateCloudInitFile(opts.CloudInit); err != nil {
		return fmt.Errorf("invalid cloud-init: %w", err)
	}
	if err := validateTimeoutValue(opts.Timeout); err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}
	return nil
}