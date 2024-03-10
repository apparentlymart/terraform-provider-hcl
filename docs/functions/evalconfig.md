# `evalconfig` function

Evaluates a HCL configuration file given as a string, and evaluates it using
a schema provided in a special schema specification file.

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
provider::hcl::evalconfig(config_src, spec_filename, variables)
```

When called, this function first parses `config_src` as HCL configuration
syntax, and then evaluates it using the schema implied by the specification
in the file at `spec_filename`.

The spec file must be written in
[this provider's specification language](../guides/spec-language.md), which
describes both the shape of the input configuration you expect and how to map
that configuration schema to a value for this function to return.

This is the most general function in this provider, but it is therefore also the
most complicated to use. If you can achieve your goal using another simpler
function -- or even not using this provider at all -- then don't use this
function.

## Example

The first step in use of this function is to describe the configuration language
you intend to accept using this provider's specification language, which is
itself based on HCL.

For this example, we'll parse and decode a contrived configuration language
shaped like the following example:

```hcl
name = "Juan"
friend {
  name = "John"
}
friend {
  name = "Yann"
}
friend {
  name = "Ermintrude"
}
```

The following specification file describes that language:

```hcl
object {
  attr "name" {
    type     = string
    required = true
  }
  default "greeting" {
    attr {
      name = "greeting"
      type = string
    }
    literal {
      value = "Hello"
    }
  }
  block_list "friends" {
    block_type = "friend"
    attr {
      name     = "name"
      type     = string
      required = true
    }
  }
}
```

If the original configuration file were in `config.hcl`, and the specification
were in the file `friends.hcldec`, then you could decode this configuration
with a function call like the following:

```hcl
provider::hcl::evalconfig(
  file("${path.module}/config.hcl"),
  "${path.module}/friends.hcldec",
  {},
)
```

The specification file describes returning the result as an object value
with the different configuration elements as attributes. In this case, the
result would be:

```hcl
{
  name     = "Juan"
  greeting = "Hello"
  friends = [
    "John",
    "Yann",
    "Ermintrude",
  ]
}
```

## Evaluating Existing Languages

This function is intended for evaluating small languages designed for local
use in a particular Terraform module.

Although it _may_ be possible to use it to decode configuration files written
in other HCL-based languages defined elsewhere, the spec language defined in
this provider is at a higher level of abstraction than most Go applications
that implement HCL-based languages, and so there are various HCL language
features that this provider does not have any means to describe.

In particular, it's not possible to evaluate the Terraform language using
this function, because Terraform includes several features that interact with
HCL at the grammar level rather than at the value level, and those behaviors
cannot be described in this provider's specification language that's designed
only to transform HCL input into dynamic values in Terraform's type system.

## Related Functions

- [`evalconfigattrs`](./evalconfigattrs.md) allows schema-free decoding of
HCL configuration files as long as the input uses only top-level attributes,
and does not include any nested blocks.
