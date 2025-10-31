# Configure HyperV Provider with SSH (Recommended)
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_port             = 22 # Default SSH port
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/id_rsa"
  # Or use password: ssh_password = "YourPassword"
}

# Alternative: WinRM configuration (legacy)
# provider "hyperv" {
#   user            = "Administrator"
#   password        = "P@ssw0rd"
#   host            = "127.0.0.1"
#   port            = 5986
#   https           = true
#   insecure        = false
#   use_ntlm        = true
#   tls_server_name = ""
#   cacert_path     = ""
#   cert_path       = ""
#   key_path        = ""
#   script_path     = "C:/Temp/terraform_%RAND%.cmd"
#   timeout         = "30s"
# }

# Create a switch
resource "hyperv_network_switch" "dmz" {
  name = "dmz_switch"
}

# Create a vhd
resource "hyperv_vhd" "webserver" {
  path = "C:/VMs/webserver.vhdx"
  size = 32 * 1024 * 1024 * 1024 # 32GB
}

# Create a machine
resource "hyperv_machine_instance" "webserver" {
  name                 = "webserver"
  generation           = 2
  processor_count      = 2
  memory_startup_bytes = 2 * 1024 * 1024 * 1024 # 2GB

  network_adaptors {
    name        = "eth0"
    switch_name = hyperv_network_switch.dmz.name
  }

  hard_disk_drives {
    path                = hyperv_vhd.webserver.path
    controller_number   = 0
    controller_location = 0
  }
}