# OpenTofu Provider Multipass Examples

This directory contains comprehensive examples demonstrating how to use the Multipass OpenTofu provider.

## Directory Structure

- `provider/` - Provider configuration examples
- `resources/` - Resource usage examples
  - `multipass_instance/` - Multipass instance resource examples
- `data-sources/` - Data source usage examples  
  - `multipass_instance/` - Multipass instance data source examples
- `complete-examples/` - Complete workflow examples
  - `vm-info-output/` - Full example that creates a VM and outputs its information

## Getting Started

1. **Install Prerequisites**:
   - Install [Multipass](https://multipass.run/) on your system
   - Ensure Go 1.23+ is installed

2. **Build and Install Provider**:
   ```bash
   make install-local
   ```

3. **Choose an Example**:
   Navigate to any example directory and follow its README instructions

## Quick Start

For a complete working example that demonstrates VM creation and information output:

```bash
cd complete-examples/vm-info-output
tofu init
tofu apply
tofu output
```

## Example Categories

### Basic Usage
- **Provider Configuration**: How to configure the provider
- **Simple Instance Creation**: Basic VM creation with default settings
- **Custom Instance Configuration**: VM with specific CPU, memory, disk settings

### Advanced Usage
- **Cloud-Init Integration**: Using cloud-init for VM customization
- **Data Sources**: Querying existing instance information
- **Import Existing Instances**: Managing pre-existing VMs with OpenTofu

### Output and Monitoring
- **Instance Information**: Retrieving VM IP addresses, state, and metadata
- **Multiple Instances**: Managing and querying multiple VMs

## Common Workflows

1. **Development Environment Setup**:
   ```bash
   cd resources/multipass_instance
   # Edit resource.tf for your needs
   tofu init && tofu apply
   ```

2. **Query Existing Infrastructure**:
   ```bash
   cd data-sources/multipass_instance
   tofu init && tofu apply
   ```

3. **Complete VM Lifecycle**:
   ```bash
   cd complete-examples/vm-info-output
   tofu init && tofu apply
   # Work with your VM
   tofu destroy
   ```