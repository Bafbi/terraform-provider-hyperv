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
  ssh                  = true
  ssh_host             = "172.16.0.200"
  ssh_user             = "nathan-admin"
  ssh_private_key_path = "~/.ssh/id_ed25519"
}

# Network switch for the test VM
resource "hyperv_network_switch" "test" {
  name = "test_switch"
}

# Upload ISO to the remote host
resource "hyperv_iso_image" "debian_installer" {
  source_iso_file_path      = "${path.module}/.images/debian-13.1.0-amd64-DVD-1.iso"
  destination_iso_file_path = "V:/.images/debian-13.1.0-amd64-DVD-1.iso"
}

# Virtual hard disk for the test VM
resource "hyperv_vhd" "test_vm_disk" {
  path = "V:/terraform_test/test_vm.vhdx"
  size = 32 * 1024 * 1024 * 1024 # 32GB
}

# Test virtual machine
resource "hyperv_machine_instance" "test_vm" {
  name                   = "terraform_test_vm"
  path                   = "V:/terraform_test"
  generation             = 2
  processor_count        = 2
  dynamic_memory         = true
  memory_startup_bytes   = 2 * 1024 * 1024 * 1024     # 2GB
  memory_minimum_bytes   = 512 * 1024 * 1024          # 512MB
  memory_maximum_bytes   = 4 * 1024 * 1024 * 1024     # 4GB
  wait_for_state_timeout = 10
  wait_for_ips_timeout   = 10
  state                  = "Off" # Keep it off for testing

  vm_firmware {
    enable_secure_boot = "off"
  }

  vm_processor {
    expose_virtualization_extensions = true
  }

  network_adaptors {
    name         = "eth0"
    switch_name  = hyperv_network_switch.test.name
    wait_for_ips = false
  }

  hard_disk_drives {
    path                = hyperv_vhd.test_vm_disk.path
    controller_number   = 0
    controller_location = 0
  }

  dvd_drives {
    controller_number   = 0
    controller_location = 1
    path                = hyperv_iso_image.debian_installer.resolve_destination_iso_file_path
  }

  depends_on = [
    hyperv_iso_image.debian_installer,
    hyperv_vhd.test_vm_disk
  ]
}

# Output the VM details
output "vm_name" {
  value = hyperv_machine_instance.test_vm.name
}

output "vm_state" {
  value = hyperv_machine_instance.test_vm.state
}

output "vhd_path" {
  value = hyperv_vhd.test_vm_disk.path
}

output "iso_path" {
  value = hyperv_iso_image.debian_installer.resolve_destination_iso_file_path
}