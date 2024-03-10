# `evalexpr` function

Evaluates a HCL expression given as a string.

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
provider::hcl::evalexpr(expr_src, variables)
```

When called, this function first parses `expr_src` as a HCL expression, and
then evaluates it in a scope that contains the variables defined in `variables`.

For example:

```
provider::hcl::evalexpr("a + b", {
  a = 1
  b = 2
})
```

The above would return the number 3, because that's the result of adding
numbers 1 and 2.

The given expression can use any of HCL's operators and can also call any of
[a small set of functions](../guides/evaluation-funcs.md).

## Related Functions

- [`evaltemplate`](./evaltemplate.md) is similar but uses HCL template syntax
instead of HCL expression syntax.
