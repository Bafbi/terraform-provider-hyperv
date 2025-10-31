# HyperV Provider Examples

This directory contains a comprehensive set of examples demonstrating how to use the Terraform HyperV provider.

## Quick Start

### SSH Connection (Recommended)

The HyperV provider supports SSH connections, which work cross-platform (Linux/macOS/Windows) and are the recommended connection method:

```hcl
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_port             = 22
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/id_rsa"
  # Or use password: ssh_password = "YourPassword"
}
```

### Path Separators

**Good news!** You can use forward slashes (`/`) in paths on Windows - no need to escape backslashes:

```hcl
resource "hyperv_vhd" "example" {
  path = "C:/VMs/my-vm.vhdx"  # ✓ Works great!
  # path = "C:\\VMs\\my-vm.vhdx"  # Also works, but more verbose
}
```

## Examples Overview

### Provider Configuration
- **[provider/](provider/)** - Basic provider configuration with SSH (recommended)
- **[provider-ssh/](provider-ssh/)** - Advanced SSH configuration options

### Complete VM Examples
- **[vm-from-scratch/](vm-from-scratch/)** - Create a VM from scratch with networking and storage
- **[clone-existing-vm/](clone-existing-vm/)** - Clone an existing VM
- **[main/](main/)** - Comprehensive example with Generation 1 and 2 VMs

### Resource Examples
- **[resources/hyperv_vhd/](resources/hyperv_vhd/)** - Virtual hard disk management
- **[resources/hyperv_machine_instance/](resources/hyperv_machine_instance/)** - Virtual machine configuration
- **[resources/hyperv_network_switch/](resources/hyperv_network_switch/)** - Network switch management
- **[resources/hyperv_iso_image/](resources/hyperv_iso_image/)** - ISO image creation

### Data Source Examples
- **[data-sources/hyperv_vhd/](data-sources/hyperv_vhd/)** - Query existing VHDs
- **[data-sources/hyperv_machine_instance/](data-sources/hyperv_machine_instance/)** - Query existing VMs
- **[data-sources/hyperv_network_switch/](data-sources/hyperv_network_switch/)** - Query existing switches

## Running Examples

Clone the repository and navigate to any example:

```bash
git clone https://github.com/Bafbi/terraform-provider-hyperv
cd terraform-provider-hyperv/examples/vm-from-scratch
terraform init
terraform plan
terraform apply
```

## Connection Methods

### SSH (Recommended)

SSH is the recommended connection method because:
- Works on Linux, macOS, and Windows
- Secure by design
- Supports key-based authentication
- Standard protocol with good tooling support

```hcl
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/id_rsa"
}
```

### WinRM (Legacy)

WinRM is still supported for backwards compatibility:

```hcl
provider "hyperv" {
  user     = "Administrator"
  password = "P@ssw0rd"
  host     = "127.0.0.1"
  port     = 5986
  https    = true
}
```

## Environment Variables

You can configure the provider using environment variables:

```bash
# SSH Configuration
export HYPERV_SSH=true
export HYPERV_SSH_HOST=hyperv-host.example.com
export HYPERV_SSH_USER=administrator
export HYPERV_SSH_PRIVATE_KEY_PATH=~/.ssh/id_rsa

# Or for password auth
export HYPERV_SSH_PASSWORD=your-password

# Then use an empty provider block
terraform {
  required_providers {
    hyperv = {
      source = "Bafbi/hyperv"
    }
  }
}

provider "hyperv" {
  # Configuration from environment variables
}
```

## Best Practices

1. **Use SSH**: Prefer SSH over WinRM for new deployments
2. **Forward Slashes**: Use forward slashes in paths - they work on Windows!
3. **Key Authentication**: Use SSH keys instead of passwords when possible
4. **Environment Variables**: Store sensitive credentials in environment variables
5. **State Management**: Use remote state for team collaboration

## Need Help?

- Check individual example README files for detailed instructions
- Review the [provider documentation](../docs/)
- Open an issue on [GitHub](https://github.com/Bafbi/terraform-provider-hyperv/issues)