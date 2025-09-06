# Basic instance
resource "multipass_instance" "example" {
  name = "my-instance"
}

# Instance with custom configuration
resource "multipass_instance" "custom" {
  name   = "custom-instance"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
}

# Instance with cloud-init
resource "multipass_instance" "with_cloud_init" {
  name       = "cloud-init-instance"
  image      = "22.04"
  cpu        = "2"
  memory     = "4G"
  disk       = "20G"
  cloud_init = "./cloud-init.yaml"

  # Configure timeouts for longer operations
  timeouts {
    create = "20m" # Cloud-init setup may take longer
    delete = "5m"
  }
}