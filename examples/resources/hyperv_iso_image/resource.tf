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

data "archive_file" "bootstrap" {
  type        = "zip"
  source_dir  = "bootstrap"
  output_path = "bootstrap.zip"
}

resource "hyperv_iso_image" "bootstrap" {
  volume_name               = "BOOTSTRAP"
  source_zip_file_path      = data.archive_file.bootstrap.output_path
  source_zip_file_path_hash = data.archive_file.bootstrap.output_sha
  destination_iso_file_path = "$env:TEMP/bootstrap.iso" # Forward slashes work in PowerShell paths!
  iso_media_type            = "dvdplusrw_duallayer"
  iso_file_system_type      = "unknown"
}