# VM Info Output Example

This example demonstrates how to create a Multipass VM instance and output its name and IP address using the locally built provider.

## Prerequisites

1. Install Multipass on your system
2. Build and install the provider locally:
   ```bash
   make install-local
   ```

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Plan the deployment:
   ```bash
   terraform plan
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

4. View the outputs:
   ```bash
   terraform output
   ```

## Expected Output

After successful deployment, you should see outputs similar to:
```
vm_name = "example-vm"
vm_ip = "192.168.64.2"
vm_state = "Running"
```

## Cleanup

To destroy the created resources:
```bash
terraform destroy
```