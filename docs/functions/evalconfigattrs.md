# `evalconfigattrs` function

Evaluates a HCL configuration file given as a string, returning a description
of any attributes it contains.

To use this function, you must first declare that your module depends on this
provider:

```hcl
terraform {
  required_providers {
    hcl = {
      source = "apparentlymart/hcl"
    }
  }
}
```

```hcl
provider::hcl::evalconfigattrs(config_src, variables)
```

When called, this function first parses `config_src` as HCL configuration
syntax, and then evaluates it in HCL's "Just Attributes" mode, which allows
dynamically detecting which attributes are declared without you needing to
provide a schema.

For example:

```
provider::hcl::evalconfigattrs(file("${path.module}/example.hcl"), {
  a = 1
  b = 2
})
```

The attributes set in the given configuration file can refer to the variables
defined in the second argument, and can call any of
[a small set of functions](../guides/evaluation-funcs.md).

"Just Attributes" mode declares that you don't expect to find any nested
blocks in the given configuration file, and so this function will return an
error if given a configuration file containing nested blocks. If you want to
design an HCL-based language that uses nested blocks then you will need to
use the more complicated function [`evalconfig`](./evalconfig.md).

If you set the second argument to the empty object `{}` then this function
effectively implements Terraform CLI's `.tfvars` file format for providing
root module input variables, but if you want to decode that format you
should use
[`provider::terraform::tfvarsdecode`](https://developer.hashicorp.com/terraform/language/functions/terraform-tfvarsdecode)
instead, to make your intentions clearer to future maintainers.

## Related Functions

- [`evalconfig`](./evalconfig.md) allows defining your own HCL-based language
using a special file format to describe your intended schema.
