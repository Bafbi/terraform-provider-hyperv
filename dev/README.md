# Local Development Directory

Test configuration for the Hyper-V provider via SSH.

## What it Creates

- Network switch: `test_switch`
- VHD: 32GB disk at `V:/terraform_test/test_vm.vhdx`
- ISO: Uploads Debian installer to remote host
- VM: 2 CPUs, 2GB RAM, connected to switch with ISO mounted

## Quick Start

```bash
# 1. Download ISO
cd .images && ./download_iso.sh && cd ..

# 2. Build and install provider
cd .. && mise run install && cd dev

# 3. Test
tofu init
tofu plan
tofu apply
```

## Configuration

Edit `test.tf` to change:
- SSH host/user/key in provider block
- VM specs (CPU, RAM, disk size)
- Paths on remote host (currently V:/terraform_test/)

## Notes

- ISO upload happens during `tofu apply` and may take a few minutes
- VM is created in "Off" state for testing
- Forward slashes work in Windows paths
- Directories are created automatically by Terraform
