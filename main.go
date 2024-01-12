package main

import (
	"context"
	"flag"
	"log"

	"github.com/pub-solar/terraform-provider-hostingde/hostingde"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name hostingde

//nolint:errcheck
func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/pub-solar/hostingde",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), hostingde.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
