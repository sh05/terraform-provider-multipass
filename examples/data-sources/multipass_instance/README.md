# Multipass Instance Data Source Examples

This directory contains examples of how to use the `multipass_instance` data source to query information about existing Multipass instances.

## Prerequisites

1. Install Multipass on your system
2. Build and install the provider locally:
   ```bash
   make install-local
   ```
3. Have at least one Multipass instance running

## Examples

### Query Specific Instance
Retrieves information about a specific instance by name:
- Instance details (IP, state, image, etc.)
- Useful for referencing existing instances

### List All Instances
Queries all available Multipass instances:
- Returns array of all instances
- Useful for inventory and monitoring

### Using Instance Data
Demonstrates how to use the retrieved data:
- Output instance IP address
- Output all instances information

## Files

- `data-source.tf` - Terraform configuration with data source examples

## Usage

1. Ensure you have Multipass instances running:
   ```bash
   multipass list
   ```

2. Update the instance name in `data-source.tf` to match your existing instance

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

4. View the outputs:
   ```bash
   terraform output
   ```

## Expected Output

The data source will return information such as:
- Instance name and state
- IPv4 addresses
- Image release information
- Resource usage (CPU load, memory, disk)
- Mount points (if any)