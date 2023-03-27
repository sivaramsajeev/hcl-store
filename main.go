package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sivaramsajeev/terraform-provider-student/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
