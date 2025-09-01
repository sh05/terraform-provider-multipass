# Terraform Provider for Multipass

A Terraform provider for managing Ubuntu virtual machines using [Canonical Multipass](https://multipass.run/).

Multipass is a lightweight VM manager for Linux, Windows and macOS that allows you to quickly create and manage Ubuntu instances.

## Features

✅ **Implemented**
- 🚀 Create, read, update, delete Multipass instances
- 📊 Query existing instances and list all instances
- ⚙️ Configure CPU, memory, disk, and Ubuntu image versions
- 🔧 Cloud-init support for instance customization
- 📋 Import existing instances into Terraform state
- 🧪 Comprehensive testing suite

🔄 **Future Enhancements** (see [Roadmap](#roadmap))
- Network configuration and port forwarding
- Volume mounting support
- Snapshot management
- Multi-platform binary releases

## Installation

### From Terraform Registry (Coming Soon)

Once published to the Terraform Registry, you can use the provider directly in your Terraform configuration:

```hcl
terraform {
  required_providers {
    multipass = {
      source  = "registry.terraform.io/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}
```

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/sh05/terraform-provider-multipass
cd terraform-provider-multipass
```

2. Build and install the provider locally:
```bash
make install-local
```

This installs the provider to `~/.terraform.d/plugins/` for local development.

## Usage

### Basic Example

```hcl
terraform {
  required_providers {
    multipass = {
      source  = "registry.terraform.io/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}

provider "multipass" {
  # Optional: specify path to multipass binary if not in PATH
  # binary_path = "/usr/local/bin/multipass"
}

# Create a basic Ubuntu instance
resource "multipass_instance" "example" {
  name   = "my-ubuntu-vm"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
}

# Get information about the instance
data "multipass_instance" "example" {
  name = multipass_instance.example.name
}

# List all instances
data "multipass_instance" "all" {}

# Output instance information
output "instance_ip" {
  value = data.multipass_instance.example.instance.ipv4
}

output "instance_state" {
  value = data.multipass_instance.example.instance.state
}
```

### Advanced Example with Cloud-Init

```hcl
resource "multipass_instance" "web_server" {
  name       = "web-server"
  image      = "22.04"
  cpu        = "2"
  memory     = "4G"
  disk       = "20G"
  cloud_init = "./cloud-init.yaml"
}
```

Create a `cloud-init.yaml` file:
```yaml
#cloud-config
packages:
  - nginx
  - git

runcmd:
  - systemctl enable nginx
  - systemctl start nginx
```

## Resources and Data Sources

### Resources

#### `multipass_instance`

Manages a Multipass Ubuntu instance.

**Arguments:**
- `name` (Required) - Instance name
- `image` (Optional) - Ubuntu image (default: latest LTS)
- `cpu` (Optional) - Number of CPUs
- `memory` (Optional) - Memory allocation (e.g., "1G", "512M")
- `disk` (Optional) - Disk space (e.g., "5G", "10G")
- `cloud_init` (Optional) - Path to cloud-init configuration file

**Attributes:**
- `id` - Instance identifier (same as name)
- `state` - Current instance state
- `ipv4` - List of IPv4 addresses assigned to the instance

### Data Sources

#### `multipass_instance`

Queries information about Multipass instances.

**Arguments:**
- `name` (Optional) - Specific instance name. If not provided, lists all instances.

**Attributes:**
- `instance` - Single instance information (when `name` is provided)
- `instances` - List of all instances (when `name` is not provided)

## Development

### Prerequisites

- Go 1.23 or later
- Terraform 1.0 or later
- Multipass installed and accessible in PATH

### Building

```bash
# Install dependencies
make deps

# Build the provider
make build

# Run tests
make test

# Run acceptance tests (requires Multipass)
make testacc

# Install locally for development
make install-local
```

### Testing

```bash
# Unit tests
go test ./...

# Acceptance tests (requires TF_ACC=1 and Multipass)
TF_ACC=1 go test ./... -v

# Test with specific instance
cd examples/complete-examples/vm-info-output
terraform init
terraform plan
terraform apply
```

## Architecture

The provider is built using the modern [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) and follows these patterns:

- **Provider**: Main provider configuration and client initialization
- **Resources**: CRUD operations for Multipass instances
- **Data Sources**: Read-only queries for existing instances
- **Client**: Go wrapper around the Multipass CLI
- **Testing**: Comprehensive unit and acceptance tests

## Roadmap

### Version 0.2.0 (Next Release)
- [ ] Network configuration support
- [ ] Port forwarding management
- [ ] Volume mounting (`multipass mount` integration)
- [ ] Better error handling and validation
- [ ] CI/CD pipeline for automated releases

### Version 0.3.0
- [ ] Snapshot management (`multipass snapshot` commands)
- [ ] Instance lifecycle operations (start, stop, suspend, restart)
- [ ] Exec operations for running commands
- [ ] Multi-platform binary releases (Linux, macOS, Windows)

### Version 1.0.0
- [ ] Production readiness
- [ ] Complete Multipass API coverage
- [ ] Enhanced documentation and examples
- [ ] Performance optimizations
- [ ] Security hardening

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Run the test suite (`make test`)
5. Update the [CHANGELOG.md](CHANGELOG.md) with your changes
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Release Process

This project uses [GoReleaser](https://goreleaser.com/) for automated releases:

```bash
# Create a new tag
git tag v0.1.0
git push origin v0.1.0

# GoReleaser will automatically create a release with binaries
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Canonical Multipass](https://multipass.run/) for the excellent VM management tool
- [HashiCorp](https://www.hashicorp.com/) for Terraform and the Plugin Framework
- The Terraform community for guidance and best practices