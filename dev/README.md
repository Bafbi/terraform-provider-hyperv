# Local Development Directory

This directory is for local development and testing of the provider.

## Quick Start

1. Build and install the provider:
   ```bash
   make install
   ```

2. Set the config file to use the local provider:
   ```bash
   export TF_CLI_CONFIG_FILE=$PWD/.tofurc
   ```

3. Initialize and test:
   ```bash
   cd dev
   tofu init
   tofu plan
   ```

## Using the Test Makefile Targets

The root Makefile includes several convenient targets for testing:

- `make test-init` - Build, install, and initialize the test directory
- `make test-plan` - Run `tofu plan` with the local provider
- `make test-apply` - Run `tofu apply` with the local provider

These targets automatically use the local provider override configuration.

## Manual Testing

You can also test manually by:

1. Building the provider: `make build`
2. Setting the config: `export TF_CLI_CONFIG_FILE=$PWD/.tofurc`
3. Working in any example directory or this `dev` directory
4. Running OpenTofu/Terraform commands as normal

The provider override will ensure your local build is used instead of downloading from the registry.
