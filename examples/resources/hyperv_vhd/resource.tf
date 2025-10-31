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

resource "hyperv_vhd" "web_server_vhd" {
  path = "C:/web_server/web_server_g2.vhdx" # Forward slashes work!
  #source               = ""
  #source_vm            = ""
  #source_disk          = 0
  vhd_type = "Dynamic"
  #parent_path          = ""
  size = 10737418240 # 10GB
  #block_size           = 0
  #logical_sector_size  = 0
  #physical_sector_size = 0
}