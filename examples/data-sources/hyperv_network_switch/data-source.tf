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

resource "hyperv_network_switch" "default" {
  name                                    = "DMZ"
  notes                                   = ""
  allow_management_os                     = true
  enable_embedded_teaming                 = false
  enable_iov                              = false
  enable_packet_direct                    = false
  minimum_bandwidth_mode                  = "None"
  switch_type                             = "Internal"
  net_adapter_names                       = []
  default_flow_minimum_bandwidth_absolute = 0
  default_flow_minimum_bandwidth_weight   = 0
  default_queue_vmmq_enabled              = false
  default_queue_vmmq_queue_pairs          = 16
  default_queue_vrss_enabled              = false
}

data "hyperv_network_switch" "default" {
  name = hyperv_network_switch.default.name
}

output "hyperv_network_switch" {
  value = data.hyperv_network_switch.default
}