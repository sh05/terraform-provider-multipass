package provider

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sh05/terraform-provider-multipass/internal/common"
)

// MultipassClient wraps the Multipass CLI
type MultipassClient struct {
	binaryPath string
}

// NewMultipassClient creates a new Multipass client
func NewMultipassClient(binaryPath string) *MultipassClient {
	if binaryPath == "" {
		binaryPath = "multipass"
	}
	return &MultipassClient{
		binaryPath: binaryPath,
	}
}

// Launch creates a new Multipass instance
func (c *MultipassClient) Launch(opts *common.LaunchOptions) error {
	args := []string{"launch"}

	if opts.Image != "" {
		args = append(args, opts.Image)
	}

	if opts.Name != "" {
		args = append(args, "--name", opts.Name)
	}

	if opts.CPU != "" {
		args = append(args, "--cpus", opts.CPU)
	}

	if opts.Memory != "" {
		args = append(args, "--memory", opts.Memory)
	}

	if opts.Disk != "" {
		args = append(args, "--disk", opts.Disk)
	}

	if opts.CloudInit != "" {
		args = append(args, "--cloud-init", opts.CloudInit)
	}

	if opts.Timeout != "" {
		// Convert duration string to seconds for multipass
		timeoutSeconds, err := c.durationToSeconds(opts.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format '%s': %w", opts.Timeout, err)
		}
		args = append(args, "--timeout", timeoutSeconds)
	}

	// Log the command being executed for debugging
	fmt.Printf("Executing multipass command: %s %s\n", c.binaryPath, strings.Join(args, " "))

	cmd := exec.Command(c.binaryPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to launch instance: %w, output: %s", err, string(output))
	}

	return nil
}

// durationToSeconds converts duration strings like "5m", "300s", "10m30s" to seconds string
func (c *MultipassClient) durationToSeconds(durationStr string) (string, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.0f", duration.Seconds()), nil
}

// GetInstance retrieves information about a specific instance
func (c *MultipassClient) GetInstance(name string) (*common.MultipassInstance, error) {
	cmd := exec.Command(c.binaryPath, "info", name, "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get instance info: %w", err)
	}

	var info common.MultipassInstanceInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse instance info: %w", err)
	}

	if len(info.Errors) > 0 {
		return nil, fmt.Errorf("multipass errors: %s", strings.Join(info.Errors, ", "))
	}

	if instance, exists := info.Info[name]; exists {
		return &instance, nil
	}

	return nil, fmt.Errorf("instance %s not found", name)
}

// ListInstances returns all instances
func (c *MultipassClient) ListInstances() ([]common.MultipassInstance, error) {
	cmd := exec.Command(c.binaryPath, "list", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var instanceList common.MultipassInstanceList
	if err := json.Unmarshal(output, &instanceList); err != nil {
		return nil, fmt.Errorf("failed to parse instance list: %w", err)
	}

	return instanceList.List, nil
}

// DeleteInstance deletes a Multipass instance
func (c *MultipassClient) DeleteInstance(name string, purge bool) error {
	// First delete the instance
	cmd := exec.Command(c.binaryPath, "delete", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w, output: %s", err, string(output))
	}

	// Then purge if requested
	if purge {
		cmd = exec.Command(c.binaryPath, "purge")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to purge instance: %w, output: %s", err, string(output))
		}
	}

	return nil
}

// StartInstance starts a stopped instance
func (c *MultipassClient) StartInstance(name string) error {
	cmd := exec.Command(c.binaryPath, "start", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start instance: %w, output: %s", err, string(output))
	}

	return nil
}

// StopInstance stops a running instance
func (c *MultipassClient) StopInstance(name string) error {
	cmd := exec.Command(c.binaryPath, "stop", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w, output: %s", err, string(output))
	}

	return nil
}

// SuspendInstance suspends a running instance
func (c *MultipassClient) SuspendInstance(name string) error {
	cmd := exec.Command(c.binaryPath, "suspend", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to suspend instance: %w, output: %s", err, string(output))
	}

	return nil
}

// RestartInstance restarts an instance
func (c *MultipassClient) RestartInstance(name string) error {
	cmd := exec.Command(c.binaryPath, "restart", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart instance: %w, output: %s", err, string(output))
	}

	return nil
}
