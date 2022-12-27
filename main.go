package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/mvantellingen/terraform-pf-bug/hashicups"
)

func main() {
	providerserver.Serve(context.Background(), hashicups.New, providerserver.ServeOpts{
		Address: "hashicorp.com/edu/hashicups-pf",
	})
}
