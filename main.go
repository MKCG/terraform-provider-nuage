package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-nuage/nuage"
)

func main() {
	tfsdk.Serve(context.Background(), nuage.New, tfsdk.ServeOpts{
 		Name: "nuage",
	})
}
