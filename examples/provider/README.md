# Provider Configuration Examples

This directory contains examples of how to configure the Multipass OpenTofu provider.

## Prerequisites

1. Install Multipass on your system
2. Build and install the provider locally:
   ```bash
   make install-local
   ```

## Provider Configuration

The Multipass provider currently requires minimal configuration. The provider automatically detects and uses the `multipass` CLI tool from your system PATH.

## Files

- `provider.tf` - Basic provider configuration example

## Usage

This provider configuration can be used as a starting point for any OpenTofu configuration using the Multipass provider:

```hcl
terraform {
  required_providers {
    multipass = {
      source = "registry.opentofu.org/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}

provider "multipass" {
  # No additional configuration required
  # Provider uses multipass CLI from system PATH
}
```

## Requirements

- Multipass must be installed and available in your system PATH
- You should be able to run `multipass version` successfully
- Sufficient permissions to create and manage Multipass instances

## Troubleshooting

If you encounter issues:

1. Verify Multipass installation:
   ```bash
   multipass version
   ```

2. Check if you can list instances:
   ```bash
   multipass list
   ```

3. Ensure the provider is properly installed:
   ```bash
   tofu init
   ```