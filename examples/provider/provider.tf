terraform {
  required_providers {
    multipass = {
      source = "sh05/multipass"
    }
  }
}

provider "multipass" {
  # Optional: specify path to multipass binary if not in PATH
  # binary_path = "/usr/local/bin/multipass"
}