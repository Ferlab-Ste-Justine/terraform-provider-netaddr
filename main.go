package main

import (
	"github.com/Ferlab-Ste-Justine/terraform-provider-netaddr/provider"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
