# Terraform HyperV Provider

Manage Microsoft Hyper-V resources via SSH from any platform.

- [Documentation](https://registry.terraform.io/providers/Bafbi/hyperv/latest/docs)
- [Issues](https://github.com/bafbi/terraform-provider-hyperv/issues)
- [Releases](https://github.com/bafbi/terraform-provider-hyperv/releases)

## Features

- Manage Hyper-V from Linux, macOS, or Windows via SSH
- Create and manage virtual machines, VHDs, and network switches
- SSH authentication with password or private key
- Full VM lifecycle: network adapters, disks, DVD drives, processors

## Requirements

- Terraform 1.0+
- Windows Server 2016+ or Windows 10+ with Hyper-V
- SSH server on Hyper-V host (OpenSSH recommended)

## Quick Start

### 1. Setup Hyper-V Host

```powershell
# Install Hyper-V
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All

# Install OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
Start-Service sshd
Set-Service -Name sshd -StartupType 'Automatic'

# Configure firewall
New-NetFirewallRule -Name sshd -DisplayName 'OpenSSH Server' `
  -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
```

### 2. Configure Provider

```hcl
provider "hyperv" {
  host     = "hyperv-host.example.com"
  user     = "Administrator"
  password = var.password
  ssh      = true
  # Or use SSH key: ssh_private_key_path = "~/.ssh/id_rsa"
}
```

### 3. Create Resources

```hcl
resource "hyperv_network_switch" "default" {
  name = "external-switch"
}

resource "hyperv_vhd" "disk" {
  path = "C:\\VMs\\disk.vhdx"
  size = 60 * 1024 * 1024 * 1024  # 60GB
}

resource "hyperv_machine_instance" "vm" {
  name         = "my-vm"
  generation   = 2
  processor_count = 2
  memory_startup_bytes = 2048 * 1024 * 1024  # 2GB
  
  network_adaptors {
    name        = "eth0"
    switch_name = hyperv_network_switch.default.name
  }
  
  hard_disk_drives {
    path = hyperv_vhd.disk.path
  }
}
```

## Provider Configuration

| Setting | Description | Default |
|---------|-------------|---------|
| `host` | Hyper-V host address | `localhost` |
| `user` | SSH username | `Administrator` |
| `password` | SSH password | - |
| `ssh_private_key_path` | Path to SSH private key | - |
| `ssh` | Enable SSH connection | `false` |
| `port` | SSH port | `22` |
| `timeout` | Connection timeout | `30s` |

Environment variables: `HYPERV_HOST`, `HYPERV_USER`, `HYPERV_PASSWORD`, `HYPERV_SSH`, etc.

## Resources

- `hyperv_network_switch` - Virtual switches
- `hyperv_vhd` - Virtual hard disks
- `hyperv_machine_instance` - Virtual machines
- `hyperv_iso_image` - ISO images

See [documentation](https://registry.terraform.io/providers/bafbi/hyperv/latest/docs) for details.

## Development

```bash
# Setup (requires mise)
mise run setup

# Build and install
mise run install

# Test
mise run test
mise run dev-plan

# Or use make
make install
make test
```

See `mise tasks` for all available commands.

## License

Mozilla Public License 2.0
