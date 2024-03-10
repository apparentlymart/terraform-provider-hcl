# HCL provider for Terraform

This provider allows a module written for Terraform, or any other
Terraform-provider-compatible software, to interact with parsing and
evaluation functionality provided by the HCL configuration language toolkit.

HCL is the toolkit used as the foundations of the Terraform language, but it's
important to realize that this provider is for _HCL itself_, and not for
the Terraform language. You can use this provider to define and evaluate
your own small HCL-based languages, but there's nothing here for evaluating
the Terraform language specifically, since that's Terraform's own
responsibility.

This provider contains only functions. Terraform v1.8 introduced the possibility
for providers to contribute functions into the language, and so this provider
effectively requires Terraform v1.8.0 or later.

## Usage

To use the functions contributed by this provider, you'll first need to import
the provider into your module as a dependency:

```hcl
terraform {
  required_providers {
    hcl = {
      source = "apparentlymart/hcl"
    }
  }
}
```

The example above uses the local name `hcl` to refer to the provider, and so
all of the functions would be called using the prefix `provider::hcl::`. For
example, the `evalexpr` function would be called as `provider::hcl::evalexpr`.
