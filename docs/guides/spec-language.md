# HCL Config Spec Language

The "spec language" is an HCL-based language used to describe other HCL-based
languages. Specifically, it describes how to project an HCL configuration with
a specific schema into a value in Terraform's type system.

This language is used only for the `evalconfig` function, which is the most
flexible evaluation function in this provider but also the most complicated.

Defining new HCL-based languages using this provider is a pretty esoteric need.
Most Terraform modules do not need the additional complexity of their own
separate HCL-based language and should instead just use Terraform language
features directly.

## Spec Format

The spec format is itself built from HCL syntax, with each HCL block serving
as a _spec_ whose block type and contents together describe a single mapping
action and, in most cases, a validation constraint. Each spec block produces
one JSON value.

A spec _file_ must have a single top-level spec block that describes the
top-level value that evaluation will return, and that spec block may have other
nested spec blocks (depending on its type) that produce nested structures and
additional validation constraints.

The most common usage is to produce an object whose attributes are derived
from the top-level content of the input file. In this case, the root of the
given spec file would have an `object` spec block whose contents describe how
each of the object's properties are to be populated using nested spec blocks.

Each spec is evaluated in the context of an HCL _body_, which is the HCL
terminology for one level of nesting in a configuration file. The top-level
objects in a file all belong to the root body of that file, and then each
nested block has its own body containing the elements within that block.
Some spec types select a new body as the context for their nested specs,
allowing nested HCL structures to be decoded.

## Spec Block Types

The following sections describe the different block types that can be used to
define specs within a spec file.

### `object` spec blocks

The `object` spec type is the most commonly used at the root of a spec file.
Its result is a JSON object whose properties are set based on any nested
spec blocks:

```hcl
object {
  attr "name" {
    type = string
  }
  block "address" {
    object {
      attr "street" {
        type = string
      }
      # ...
    }
  }
}
```

Nested spec blocks inside `object` must always have an extra block label
(`"name"`, `"address"` and `"street"` in the above example) that specifies
the name of the attribute that should be created in the resulting object.

This label also acts as a default name selector for the nested spec, allowing
the `attr` blocks in the above example to omit the usually-required `name`
argument in cases where the HCL input name and JSON output name are the same.

An `object` spec block creates no validation constraints, but it passes on
any validation constraints created by the nested specs.

### `tuple` spec blocks

The `tuple` spec type produces a tuple whose elements are set based on any
nested spec blocks:

```hcl
tuple {
  attr {
    name = "first_element"
    type = string
  }
  attr {
    name = "second_element"
    type = string
  }
}
```

A `tuple` spec block creates no validation constraints, but it passes on
any validation constraints created by the nested specs.

### `attr` spec blocks

The `attr` spec type reads the value of an attribute in the current body
and returns that value as its result. It also creates validation constraints
for the given attribute name and its value.

```hcl
attr {
  name     = "document_root"
  type     = string
  required = true
}
```

`attr` spec blocks accept the following arguments:

* `name` (required) - The attribute name to expect within the HCL input file.
  This may be omitted when a default name selector is created by a parent
  `object` spec, if the input attribute name should match the output JSON
  object property name.

* `type` (optional) - A [type expression](#type-expressions) that the given
  attribute value must conform to. If this argument is set, the provider will
  automatically convert the given input value to this type or produce an
  error if that is not possible.

* `required` (optional) - If set to `true`, the provider will produce an error
  if a value is not provided for the source attribute.

`attr` is a leaf spec type, so no nested spec blocks are permitted.

### `block` spec blocks

The `block` spec type applies one nested spec block to the contents of a
block within the current body and returns the result of that spec. It also
creates validation constraints for the given block type name.

```hcl
block {
  block_type = "logging"

  object {
    attr "level" {
      type = string
    }
    attr "file" {
      type = string
    }
  }
}
```

`block` spec blocks accept the following arguments:

* `block_type` (required) - The block type name to expect within the HCL
  input file. This may be omitted when a default name selector is created
  by a parent `object` spec, if the input block type name should match the
  output attribute name.

* `required` (optional) - If set to `true`, the provider will produce an error
  if a block of the specified type is not present in the current body.

`block` creates a validation constraint that there must be zero or one blocks
of the given type name, or exactly one if `required` is set.

`block` expects a single nested spec block, which is applied to the body of
the block of the given type when it is present.

### `block_list` spec blocks

The `block_list` spec type is similar to `block`, but it accepts zero or
more blocks of a specified type rather than requiring zero or one. The
result is a tuple with one entry per block of the given type.

```hcl
block_list {
  block_type = "log_file"

  object {
    attr "level" {
      type = string
    }
    attr "filename" {
      type     = string
      required = true
    }
  }
}
```

`block_list` spec blocks accept the following arguments:

* `block_type` (required) - The block type name to expect within the HCL
  input file. This may be omitted when a default name selector is created
  by a parent `object` spec, if the input block type name should match the
  output JSON object property name.

* `min_items` (optional) - If set to a number greater than zero, the provider
  will produce an error if fewer than the given number of blocks are present.

* `max_items` (optional) - If set to a number greater than zero, the provider
  will produce an error if more than the given number of blocks are present.
  This attribute must be greater than or equal to `min_items` if both are set.

`block` creates a validation constraint on the number of blocks of the given
type that must be present.

`block` expects a single nested spec block, which is applied to the body of
each matching block to produce the resulting list items.

### `block_set` spec blocks

The `block_set` spec type behaves the same as `block_list` except that
the result is in no specific order and any duplicate items are removed.

```hcl
block_set {
  block_type = "log_file"

  object {
    attr "level" {
      type = string
    }
    attr "filename" {
      type     = string
      required = true
    }
  }
}
```

The contents of `block_set` are the same as for `block_list`.

### `block_map` spec blocks

The `block_map` spec type is similar to `block`, but it accepts zero or
more blocks of a specified type rather than requiring zero or one. The
result is an object, or possibly multiple nested objects, whose properties
are derived from the labels set on each matching block.

```hcl
block_map {
  block_type = "log_file"
  labels = ["filename"]

  object {
    attr "level" {
      type     = string
      required = true
    }
  }
}
```

`block_map` spec blocks accept the following arguments:

* `block_type` (required) - The block type name to expect within the HCL
  input file. This may be omitted when a default name selector is created
  by a parent `object` spec, if the input block type name should match the
  output attribute name.

* `labels` (required) - A list of user-oriented block label names. Each entry
  in this list creates one level of object within the output value, and
  requires one additional block header label on any child block of this type.
  Block header labels are the quoted strings that appear after the block type
  name but before the opening `{`.

`block` creates a validation constraint on the number of labels that blocks
of the given type must have.

`block` expects a single nested spec block, which is applied to the body of
each matching block to produce the resulting map items.

## `block_attrs` spec blocks

The `block_attrs` spec type is similar to an `attr` spec block of a map type,
but it produces an object from all the attributes of a block together, rather
than from an individual attribute's expression.

```hcl
block_attrs {
  block_type   = "variables"
  element_type = string
  required     = false
}
```

This allows a map with user-defined keys to be produced within block syntax,
but due to the constraints of that syntax it also means that the user will
be unable to dynamically-generate either individual key names using key
expressions or the entire map value using a `for` expression.

`block_attrs` spec blocks accept the following arguments:

* `block_type` (required) - The block type name to expect within the HCL
  input file. This may be omitted when a default name selector is created
  by a parent `object` spec, if the input block type name should match the
  output attribute name.

* `element_type` (required) - The value type to require for each of the
  attributes within a matched block. The resulting value will be an
  object whose attribute values are all of this type.

* `required` (optional) - If `true`, an error will be produced if a block
  of the given type is not present. If `false` -- the default -- an absent
  block will be indicated by producing `null`.

## `literal` spec blocks

The `literal` spec type returns a given literal value, and creates no
validation constraints. It is most commonly used with the `default` spec
type to create a fallback value, but can also be used e.g. to fill out
required properties in an `object` spec that do not correspond to any
construct in the input configuration.

```hcl
literal {
  value = "hello world"
}
```

`literal` spec blocks accept the following argument:

* `value` (required) - The value to return. This attribute may be an expression
  that uses [evaluation functions](./evaluation-funcs.md).

`literal` is a leaf spec type, so no nested spec blocks are permitted.

## `default` spec blocks

The `default` spec type evaluates a sequence of nested specs in turn and
returns the result of the first one that produces a non-null value.
It creates no validation constraints of its own, but passes on the validation
constraints from its first nested block.

```hcl
default {
  attr {
    name = "private"
    type = bool
  }
  literal {
    value = false
  }
}
```

A `default` spec block must have at least one nested spec block, and should
generally have at least two since otherwise the `default` wrapper is a no-op.

The second and any subsequent spec blocks are _fallback_ specs. These exhibit
their usual behavior but are not able to impose validation constraints on the
current body since they are not evaluated unless all prior specs produce
`null` as their result.

## `transform` spec blocks

The `transform` spec type evaluates one nested spec and then evaluates a given
expression with that nested spec result to produce a final value.
It creates no validation constraints of its own, but passes on the validation
constraints from its nested block.

```hcl
transform {
  attr {
    name = "size_in_mb"
    type = number
  }

  # Convert result to a size in bytes
  result = nested * 1024 * 1024
}
```

`transform` spec blocks accept the following argument:

* `result` (required) - The expression to evaluate on the result of the nested
  spec. The variable `nested` is defined when evaluating this expression, with
  the result value of the nested spec.

The `result` expression may use [evaluation functions](./evaluation-funcs.md).

## Predefined Variables

The `evalconfig` function accepts values for variables to expose into the input
file's expression scope in its third argument, and this is the most common way
to pass values since it allows them to be dynamically populated by the calling
module.

However, it's also possible to pre-define variables with constant values
within a spec file, using the top-level `variables` block type:

```hcl
variables {
  name = "Stephen"
}
```

Variables of the same name defined via the `evalconfig` argument will override
predefined variables of the same name, so this mechanism can also be used to
provide defaults for variables that are overridden only in certain contexts.

## Custom Functions

The spec can make arbitrary HCL functions available in the input file's
expression scope, and thus allow simple computation within the input file,
in addition to HCL's built-in operators.

Custom functions are defined in the spec file with the top-level `function`
block type:

```
function "add_one" {
  params = [n]
  result = n + 1
}
```

Functions behave in a similar way to the `transform` spec type in that the
given `result` attribute expression is evaluated with additional variables
defined with the same names as the defined `params`.

The [evaluation functions](./evaluation-funcs.md) can be used within
custom function expressions, allowing them to be optionally exposed into the
input file:

```
function "upper" {
  params = [str]
  result = upper(str)
}

function "min" {
  params         = []
  variadic_param = nums
  result         = min(nums...)
}
```

Custom functions defined in the spec cannot be called from the spec itself.

## Type Expressions

Type expressions are used to describe the expected type of an attribute, as
an additional validation constraint.

A type expression uses primitive type names and compound type constructors.
A type constructor builds a new type based on one or more type expression
arguments.

The following type names and type constructors are supported:

* `any` is a wildcard that accepts a value of any type.
* `string` is a Unicode string.
* `number` is an arbitrary-precision floating point number.
* `bool` is a boolean value (`true` or `false`)
* `list(element_type)` constructs a list type with the given element type
* `set(element_type)` constructs a set type with the given element type
* `map(element_type)` constructs a map type with the given element type
* `object({name1 = element_type, name2 = element_type, ...})` constructs
  an object type with the given attribute types.
* `tuple([element_type, element_type, ...])` constructs a tuple type with
  the given element types. This can be used, for example, to require an
  array with a particular number of elements, or with elements of different
  types.

`null` is a valid value of any type, and not a type itself.
