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
  ssh_port             = 22
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/id_rsa"
}

data "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}

data "hyperv_machine_instance" "web_server_g1" {
  name = "web_server_g1"
}

resource "hyperv_vhd" "web_server_g3_vhd" {
  path      = "C:/VhdX2/web_Server_g3.vhdx" # Forward slashes work!
  source_vm = data.hyperv_machine_instance.web_server_g1.name
}

resource "hyperv_machine_instance" "web_server_g3" {
  name          = "web_server_g3"
  static_memory = true

  network_adaptors {
    name        = "wan"
    switch_name = data.hyperv_network_switch.dmz_network_switch.name
  }

  hard_disk_drives {
    path                = hyperv_vhd.web_server_g3_vhd.path
    controller_number   = "0"
    controller_location = "0"
  }

  integration_services = {
    VSS = true
  }
}