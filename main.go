package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/flexkube/terraform-provider-flexkube/flexkube"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexkube.Provider,
	})
}
