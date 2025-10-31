terraform {
  required_providers {
    hyperv = {
      source  = "Bafbi/hyperv"
      version = ">= 1.3.0"
    }
  }
}

# SSH connection (Recommended)
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/id_rsa"
}

data "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}
