terraform {
  required_providers {
    multipass = {
      source  = "registry.opentofu.org/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}

provider "multipass" {
  # Optional: specify path to multipass binary if not in PATH
  # binary_path = "/usr/local/bin/multipass"
}