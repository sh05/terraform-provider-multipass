# Multipass Instance Resource Examples

This directory contains examples of how to use the `multipass_instance` resource.

## Prerequisites

1. Install Multipass on your system
2. Build and install the provider locally:
   ```bash
   make install-local
   ```

## Examples

### Basic Instance
Creates a simple Multipass instance with default settings:
- Default image (latest Ubuntu LTS)
- Default CPU, memory, and disk allocations

### Custom Instance
Creates an instance with custom specifications:
- Ubuntu 22.04 image
- 2 CPUs
- 2GB memory
- 10GB disk

### Instance with Cloud-Init
Creates an instance with cloud-init configuration:
- Custom hardware specifications
- Cloud-init script from `cloud-init.yaml`

## Files

- `resource.tf` - Terraform configuration with all three examples
- `cloud-init.yaml` - Cloud-init configuration for the third example
- `import.sh` - Example script for importing existing instances

## Usage

1. Choose which example to use and modify `resource.tf` accordingly
2. Initialize Terraform:
   ```bash
   terraform init
   ```
3. Plan and apply:
   ```bash
   terraform plan
   terraform apply
   ```

## Verifying Cloud-Init Configuration

For the cloud-init example, you can verify that the configuration was applied correctly:

1. **Check instance status**:
   ```bash
   multipass info cloud-init-instance
   ```

2. **Verify nginx is running**:
   ```bash
   multipass exec cloud-init-instance -- systemctl status nginx
   ```

3. **Check the custom web page**:
   ```bash
   multipass exec cloud-init-instance -- cat /var/www/html/index.html
   ```

4. **Test HTTP access**:
   ```bash
   # Get the instance IP
   INSTANCE_IP=$(multipass info cloud-init-instance --format json | grep -o '"ipv4":\["[^"]*"' | cut -d'"' -f4)
   
   # Test with curl
   curl http://$INSTANCE_IP
   ```

5. **Check cloud-init logs**:
   ```bash
   multipass exec cloud-init-instance -- cloud-init status
   multipass exec cloud-init-instance -- tail -20 /var/log/cloud-init-output.log
   ```

## Import Existing Instances

To import an existing Multipass instance, use the import script:
```bash
./import.sh
```

Or manually:
```bash
terraform import multipass_instance.example instance-name
```