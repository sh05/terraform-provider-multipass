# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of Terraform provider for Multipass
- `multipass_instance` resource for managing VM instances
- `multipass_instance` data source for querying instance information
- Support for custom CPU, memory, and disk configurations
- Cloud-init integration for VM customization
- Import support for existing instances
- Comprehensive examples and documentation

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [0.1.0] - TBD

### Added
- Initial release of terraform-provider-multipass
- Basic CRUD operations for Multipass instances
- Data source for querying instance information
- Examples and documentation

[Unreleased]: https://github.com/sh05/terraform-provider-multipass/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/sh05/terraform-provider-multipass/releases/tag/v0.1.0