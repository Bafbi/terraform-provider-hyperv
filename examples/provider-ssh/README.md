# HyperV Provider - SSH Configuration Examples

This directory contains examples of how to configure the HyperV provider to use SSH instead of WinRM.

## Prerequisites

### Windows Host with SSH

1. **Install OpenSSH Server** (Windows 10 1809+/Server 2019+):
   ```powershell
   Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
   Start-Service sshd
   Set-Service -Name sshd -StartupType 'Automatic'
   ```

2. **Install PowerShell Core**:
   ```powershell
   winget install Microsoft.PowerShell
   ```

3. **Set PowerShell as default SSH shell** (optional):
   ```powershell
   New-ItemProperty -Path "HKLM:\SOFTWARE\OpenSSH" -Name DefaultShell `
     -Value "C:\Program Files\PowerShell\7\pwsh.exe" -PropertyType String -Force
   ```

### Linux Host with PowerShell

1. Install PowerShell Core
2. Install HyperV management tools (if available)
3. Ensure SSH server is running

## Authentication Methods

### 1. Password Authentication

```hcl
provider "hyperv" {
  ssh          = true
  ssh_host     = "hyperv-host.example.com"
  ssh_user     = "administrator"
  ssh_password = "your-password"
}
```

**Using environment variables** (recommended for passwords):
```bash
export HYPERV_SSH=true
export HYPERV_SSH_HOST=hyperv-host.example.com
export HYPERV_SSH_USER=administrator
export HYPERV_SSH_PASSWORD=your-password
```

### 2. Private Key Authentication (File Path)

```hcl
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/hyperv_rsa"
}
```

**Using environment variables**:
```bash
export HYPERV_SSH=true
export HYPERV_SSH_HOST=hyperv-host.example.com
export HYPERV_SSH_USER=administrator
export HYPERV_SSH_PRIVATE_KEY_PATH=~/.ssh/hyperv_rsa
```

### 3. Private Key Authentication (Content)

```hcl
provider "hyperv" {
  ssh             = true
  ssh_host        = "hyperv-host.example.com"
  ssh_user        = "administrator"
  ssh_private_key = file("~/.ssh/hyperv_rsa")
}
```

Or using a variable:
```hcl
variable "ssh_private_key" {
  type      = string
  sensitive = true
}

provider "hyperv" {
  ssh             = true
  ssh_host        = "hyperv-host.example.com"
  ssh_user        = "administrator"
  ssh_private_key = var.ssh_private_key
}
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `ssh` | bool | `false` | Enable SSH instead of WinRM |
| `ssh_host` | string | value of `host` | SSH server hostname or IP |
| `ssh_port` | int | `22` | SSH server port |
| `ssh_user` | string | value of `user` | SSH username |
| `ssh_password` | string | value of `password` | SSH password (sensitive) |
| `ssh_private_key` | string | `""` | SSH private key content (sensitive) |
| `ssh_private_key_path` | string | `""` | Path to SSH private key file |

## Environment Variables

All configuration options can be set via environment variables:

- `HYPERV_SSH` - Enable SSH (true/false)
- `HYPERV_SSH_HOST` - SSH host
- `HYPERV_SSH_PORT` - SSH port
- `HYPERV_SSH_USER` - SSH username
- `HYPERV_SSH_PASSWORD` - SSH password
- `HYPERV_SSH_PRIVATE_KEY` - SSH private key content
- `HYPERV_SSH_PRIVATE_KEY_PATH` - Path to SSH private key

## SSH Key Setup

### Generate SSH Key Pair

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/hyperv_rsa -C "hyperv-terraform"
```

### Copy Public Key to Windows Host

#### Option 1: Using ssh-copy-id (if available)
```bash
ssh-copy-id -i ~/.ssh/hyperv_rsa.pub administrator@hyperv-host
```

#### Option 2: Manual copy
```bash
cat ~/.ssh/hyperv_rsa.pub | ssh administrator@hyperv-host "mkdir -p .ssh && cat >> .ssh/authorized_keys"
```

#### Option 3: PowerShell on Windows
```powershell
# On Windows host
$key = Get-Content C:\path\to\public\key.pub
Add-Content -Path $env:USERPROFILE\.ssh\authorized_keys -Value $key
```

### Set Correct Permissions

On Linux/Mac client:
```bash
chmod 600 ~/.ssh/hyperv_rsa
chmod 644 ~/.ssh/hyperv_rsa.pub
```

On Windows host:
```powershell
icacls "$env:USERPROFILE\.ssh\authorized_keys" /inheritance:r
icacls "$env:USERPROFILE\.ssh\authorized_keys" /grant:r "$env:USERNAME:R"
```

## Testing SSH Connection

Before using with Terraform, test the SSH connection:

```bash
# Test basic connection
ssh -i ~/.ssh/hyperv_rsa administrator@hyperv-host

# Test PowerShell availability
ssh -i ~/.ssh/hyperv_rsa administrator@hyperv-host "pwsh -Command 'Get-Host'"

# Test HyperV cmdlet availability
ssh -i ~/.ssh/hyperv_rsa administrator@hyperv-host "pwsh -Command 'Get-Command Get-VM'"
```

## Fallback Configuration

If you don't specify SSH-specific options, the provider will fall back to general options:

```hcl
provider "hyperv" {
  ssh      = true
  host     = "hyperv-host.example.com"  # Used as ssh_host
  user     = "administrator"             # Used as ssh_user
  password = "your-password"             # Used as ssh_password
}
```

This is useful for easy migration from WinRM to SSH.

## Example: Complete Configuration

```hcl
terraform {
  required_providers {
    hyperv = {
      source = "registry.terraform.io/taliesins/hyperv"
    }
  }
}

provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_port             = 22
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/hyperv_rsa"
  timeout              = "30s"
}

resource "hyperv_network_switch" "my_switch" {
  name = "my_virtual_switch"
}

resource "hyperv_machine_instance" "my_vm" {
  name             = "my-test-vm"
  generation       = 2
  processor_count  = 2
  memory_startup_bytes = 2147483648  # 2GB
  
  network_adaptors {
    name        = "eth0"
    switch_name = hyperv_network_switch.my_switch.name
  }
}
```

## Troubleshooting

### Connection refused
- Check SSH service is running: `Get-Service sshd`
- Check firewall allows port 22
- Verify host and port are correct

### Permission denied (publickey)
- Ensure public key is in `~/.ssh/authorized_keys` on host
- Check file permissions on both client and host
- Try with password authentication first to verify other settings

### PowerShell command not found
- Install PowerShell Core on the host
- Verify `pwsh` is in PATH: `ssh user@host "which pwsh"`

### SSH works but Terraform fails
- Enable debug logging: `export TF_LOG=DEBUG`
- Check the log for detailed error messages
- Verify HyperV PowerShell modules are available

## Migration from WinRM

To migrate from WinRM to SSH, simply add `ssh = true` and SSH credentials:

### Before (WinRM)
```hcl
provider "hyperv" {
  host     = "hyperv-host.example.com"
  port     = 5986
  user     = "administrator"
  password = "your-password"
  https    = true
  insecure = true
}
```

### After (SSH)
```hcl
provider "hyperv" {
  ssh      = true
  host     = "hyperv-host.example.com"  # Reused as ssh_host
  user     = "administrator"             # Reused as ssh_user
  password = "your-password"             # Reused as ssh_password
}
```

Or use SSH-specific fields:
```hcl
provider "hyperv" {
  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/hyperv_rsa"
}
```

## Security Best Practices

1. **Use private key authentication** instead of passwords
2. **Use environment variables** for sensitive values
3. **Restrict SSH access** with firewall rules
4. **Use strong keys**: RSA 4096-bit or Ed25519
5. **Rotate credentials** regularly
6. **Enable SSH audit logging** on the host
7. **Use dedicated service accounts** for Terraform

## Additional Resources

- [OpenSSH Server Installation (Windows)](https://docs.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse)
- [PowerShell Core Installation](https://docs.microsoft.com/en-us/powershell/scripting/install/installing-powershell)
- [SSH Key-Based Authentication](https://www.ssh.com/academy/ssh/public-key-authentication)
