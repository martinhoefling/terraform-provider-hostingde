package main

import (
	"context"

	"github.com/pub-solar/terraform-provider-hostingde/hostingde"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), hostingde.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/hostingde/hostingde",
	})
}
