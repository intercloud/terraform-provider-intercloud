package main

import (
	"github.com/intercloud/terraform-provider-intercloud/intercloud"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: intercloud.Provider,
	})
}
