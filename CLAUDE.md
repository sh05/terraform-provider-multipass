# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **completed and functional** OpenTofu provider for Canonical Multipass - a tool for creating Ubuntu VMs on demand. The provider is built using the modern Terraform Plugin Framework and provides full CRUD operations for managing Multipass instances.

## Implemented Architecture

The provider follows standard Terraform provider patterns with this structure:

- **Provider Registration**: Main provider configuration and schema
- **Resources**: Individual Terraform resources (instances, networks, etc.)
- **Data Sources**: Read-only data sources for querying Multipass state
- **Client Layer**: Go client for interacting with Multipass CLI/API
- **Documentation**: Provider and resource documentation

## Common Development Commands

This project uses standard Go and OpenTofu provider tooling:

```bash
# Build the provider
go build -o opentofu-provider-multipass

# Run tests
go test ./...

# Generate documentation (using terraform-plugin-docs)
go generate

# Install provider locally for testing
make install-local
```

## OpenTofu Provider Development Notes

- OpenTofu providers are written in Go using the Terraform Plugin SDK (compatible)
- Resources should implement CRUD operations (Create, Read, Update, Delete)
- Provider configuration typically includes authentication details for the target system
- Integration tests often require the actual service (Multipass) to be installed
- Documentation is usually auto-generated from schema descriptions

## Multipass Integration

- Multipass CLI commands are typically used via `multipass` binary
- Common operations: launch, delete, list, info, mount, exec
- Instance management includes CPU, memory, disk configuration
- Network configuration and port forwarding capabilities

