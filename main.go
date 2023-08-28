// Package main serves as a main entrypoint for the terraform-provider-flexkube.
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/flexkube/terraform-provider-flexkube/flexkube"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexkube.Provider,
	})
}
