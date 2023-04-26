package main

import (
	"context"

	"github.com/pub-solar/terraform-provider-hostingde/hostingde"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name hostingde

func main() {
	providerserver.Serve(context.Background(), hostingde.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pub-solar/hostingde",
	})
}
