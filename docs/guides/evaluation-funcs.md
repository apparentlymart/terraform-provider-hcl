# Evaluation Functions

This provider defines various functions that evaluate expressions written in
HCL expression syntax, either directly or indirectly.

Those expressions can therefore contain their own function calls. The
dynamically-evaluated expressions have access to the following set of functions,
each of which is defined similarly to a function in Terraform itself.

| Function | Description |
|--|--|
| `abs(number)` | Discard the sign of the given number. |
| `coalesce(any...)` | Return the first argument whose value isn't `null`. |
| `concat(seqs...)` | Concatenate together all of the given lists or tuples and return the result. |
| `hasindex(collection, index)` | Returns true if `collection[index]` would succeed. |
| `int(number)` | Returns the integer component of the given number, rounding toward zero. |
| `jsondecode(string)` | Evaluates the given string as a JSON document and returns a value representing the result. |
| `jsonencode(any)` | Returns a JSON representation of the given value. |
| `length(collection)` | Returns the number of elements in the given collection. |
| `lower(string)` | Returns the given string with all of the caseable letters converted to lowercase. |
| `max(numbers...)` | Returns the greatest of all of the given numbers. |
| `min(numbers...)` | Returns the least of all of the given numbers. |
| `reverse(seq)` | Returns the given sequence with its elements in reverse order. |
| `strlen(string)` | Returns the number of characters (grapheme clusters) in the given string. |
| `substr(string, start, length)` | Returns a substring of the given string. |
| `upper(string)` | Returns the given string with all of the caseable letters converted to uppercase. |

Note that these are not functions exported by this provider for use in a module
that uses this provider. Instead, they are available only indirectly through
this provider's evaluation functions, like `evalexpr`.

(A number of Terraform's built-in functions are really provided by a third-party
library called `cty`, which is also used by this provider. All of the functions
above belong to `cty` rather than to Terraform, and so may not behave exactly
the same as the similar Terraform functions, since Terraform has its own
customized implementations of some functions.)
