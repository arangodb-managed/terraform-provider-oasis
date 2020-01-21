package main

import (
	"github.com/arangodb-managed/terraform-provider-oasis/oasis"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return oasis.Provider()
		},
	})
}
