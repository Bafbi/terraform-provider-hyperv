provider_installation {
  dev_overrides {
    "Bafbi/hyperv" = "__GOBIN__"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal.
  direct {}
}
