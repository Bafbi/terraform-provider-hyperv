# Example HyperV provider configuration using SSH

# Option 1: SSH with password authentication
provider "hyperv" {
  ssh          = true
  ssh_host     = "hyperv-host.example.com"
  ssh_port     = 22
  ssh_user     = "administrator"
  ssh_password = "your-password" # Better to use environment variable
}

# Option 2: SSH with private key from file
provider "hyperv" {
  alias = "ssh_key"

  ssh                  = true
  ssh_host             = "hyperv-host.example.com"
  ssh_port             = 22
  ssh_user             = "administrator"
  ssh_private_key_path = "~/.ssh/hyperv_rsa"
}

# Option 3: SSH with private key content
provider "hyperv" {
  alias = "ssh_key_content"

  ssh             = true
  ssh_host        = "hyperv-host.example.com"
  ssh_port        = 22
  ssh_user        = "administrator"
  ssh_private_key = file("~/.ssh/hyperv_rsa") # Read from file
}

# Option 4: Using environment variables (recommended)
# Set these environment variables:
#   export HYPERV_SSH=true
#   export HYPERV_SSH_HOST=hyperv-host.example.com
#   export HYPERV_SSH_USER=administrator
#   export HYPERV_SSH_PASSWORD=your-password
# or
#   export HYPERV_SSH_PRIVATE_KEY_PATH=~/.ssh/hyperv_rsa

provider "hyperv" {
  alias = "from_env"
  # All configuration comes from environment variables
}

# Option 5: Fallback to general configuration
# If ssh_host, ssh_user, etc. are not specified,
# the provider will use host, user, password from general config
provider "hyperv" {
  alias = "fallback"

  ssh      = true
  host     = "hyperv-host.example.com" # Used as ssh_host
  user     = "administrator"           # Used as ssh_user
  password = var.password              # Used as ssh_password
}
