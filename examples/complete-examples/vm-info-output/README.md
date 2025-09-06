# VM Info Output Example

This example demonstrates how to create a Multipass VM instance and output its name and IP address using the locally built provider.

## Prerequisites

1. Install Multipass on your system
2. Build and install the provider locally:
   ```bash
   make install-local
   ```

## Usage

1. Initialize OpenTofu:
   ```bash
   tofu init
   ```

2. Plan the deployment:
   ```bash
   tofu plan
   ```

3. Apply the configuration:
   ```bash
   tofu apply
   ```

4. View the outputs:
   ```bash
   tofu output
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
tofu destroy
```