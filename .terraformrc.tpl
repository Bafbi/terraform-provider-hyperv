provider_installation {
  dev_overrides {
    "registry.terraform.io/taliesins/hyperv" = "__GOBIN__"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal.
  direct {}
}
