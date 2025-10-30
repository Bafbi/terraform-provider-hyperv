# Local development test configuration
# This is a minimal test file for quick provider testing

terraform {
  required_providers {
    hyperv = {
      source = "registry.terraform.io/taliesins/hyperv"
    }
  }
}

provider "hyperv" {
  ssh = true
  ssh_host = "172.16.0.200"
  ssh_user = "nathan-admin"
  ssh_private_key_path = "~/.ssh/id_ed25519"
}

resource "hyperv_network_switch" "test" {
    name = "test_switch"
}