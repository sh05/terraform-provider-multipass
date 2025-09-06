terraform {
  required_providers {
    multipass = {
      source = "registry.terraform.io/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}

provider "multipass" {}

# Create a VM instance
resource "multipass_instance" "example" {
  name   = "example-vm"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
}

# Get VM information using data source
data "multipass_instance" "example_info" {
  name = multipass_instance.example.name
  depends_on = [multipass_instance.example]
}

# Output VM name and IP
output "vm_name" {
  description = "Name of the created VM"
  value       = multipass_instance.example.name
}

output "vm_ip" {
  description = "IP address of the created VM"
  value       = length(data.multipass_instance.example_info.instance.ipv4) > 0 ? data.multipass_instance.example_info.instance.ipv4[0] : "No IP assigned"
}

output "vm_state" {
  description = "Current state of the VM"
  value       = data.multipass_instance.example_info.instance.state
}