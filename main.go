// terraform-provider-hcl exposes the parser, decoder, and expression evaluation
// functionality of HCL as functions for use in Terraform.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apparentlymart/terraform-provider-hcl/internal/tfhcl"
)

func main() {
	provider := tfhcl.NewProvider()
	err := provider.Serve(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start provider: %s", err)
		os.Exit(1)
	}
}
