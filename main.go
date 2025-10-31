package main

import (
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/taliesins/terraform-provider-hyperv/internal/provider"
)

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "0.0.0"

	// goreleaser can also pass the specific commit if you want
	commit string = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Remove duplicated timestamp from logs to make them more readable (see: https://developer.hashicorp.com/terraform/plugin/log/writing#legacy-log-troubleshooting)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	opts := &plugin.ServeOpts{
		Debug: debugMode,

		ProviderAddr: "registry.terraform.io/bafbi/hyperv",

		ProviderFunc: provider.New(version, commit),
	}

	plugin.Serve(opts)
}
