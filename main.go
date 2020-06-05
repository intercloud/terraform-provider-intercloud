package main

import (
	"gitlab.intercloud.fr/intercloud/ds9/iac/terraform-provider-intercloud.git/intercloud"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: intercloud.Provider,
	})
}
