# Get information about a specific instance
data "multipass_instance" "example" {
  name = "my-instance"
}

# List all instances
data "multipass_instance" "all" {
}

# Use instance data
output "instance_ip" {
  value = data.multipass_instance.example.instance.ipv4
}

output "all_instances" {
  value = data.multipass_instance.all.instances
}