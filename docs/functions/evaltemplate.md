# `evaltemplate` function

Evaluates a HCL template given as a string.

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
provider::hcl::evaltemplate(template_src, variables)
```

When called, this function first parses `template_src` as a HCL string template,
and then evaluates it in a scope that contains the variables defined in
`variables`.

For example:

```
provider::hcl::evaltemplate("The sum is $${a + b}", {
  a = 1
  b = 2
})
```

The above would return the string `"The sum is 3"`, because that's the result
of adding numbers 1 and 2.

It's important to escape any interpolation or control sequences included in
your template if you're writing the template as a quoted or heredoc string
directly inside your module, because otherwise Terraform will attempt to
evaluate the sequences itself before calling the function. For larger templates
or templates that include lots of interpolation or control sequences, consider
placing the template in a separate file and loading it using Terraform's
built-in function `file`:

```
provider::hcl::evaltemplate(file("${path.module}/large-template.tmpl"), {
  a = 1
  b = 2
})
```

The given template can use any of HCL's operators and can also call any of
[a small set of functions](../guides/evaluation-funcs.md).

## Related Functions

- [`evalexpr`](./evalexpr.md) is similar but uses HCL expression syntax
instead of HCL template syntax.
