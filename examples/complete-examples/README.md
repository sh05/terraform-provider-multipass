# Complete Examples

This directory contains complete, end-to-end examples that demonstrate real-world usage scenarios of the Multipass Terraform provider.

## Available Examples

### VM Info Output (`vm-info-output/`)
A comprehensive example that demonstrates:
- Creating a Multipass VM instance
- Using data sources to query instance information  
- Outputting VM name, IP address, and state
- Complete Terraform workflow from creation to cleanup

## Usage Pattern

Each complete example follows this pattern:

1. **Complete Configuration**: Self-contained Terraform configuration
2. **Real-world Scenario**: Addresses actual use cases
3. **Documented Workflow**: Step-by-step instructions
4. **Verification Steps**: How to confirm everything works
5. **Cleanup Instructions**: How to tear down resources

## Getting Started

1. Choose an example directory
2. Follow the README instructions in that directory
3. Ensure you have built and installed the provider locally:
   ```bash
   make install-local
   ```

## Difference from Basic Examples

- **Basic Examples** (`../resources/`, `../data-sources/`): Focus on individual resource or data source usage
- **Complete Examples**: Show full workflows combining multiple components

These examples are ideal for:
- Learning complete Terraform workflows
- Testing provider functionality
- Starting point for your own configurations
- CI/CD pipeline testing