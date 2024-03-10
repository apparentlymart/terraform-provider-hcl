package tfhcl

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var evalconfigattrsFunc = &function.Spec{
	Description: "Evaluates HCL configuration source code in \"Just Attributes\" mode, returning an object representing the attributes.",
	Params: []function.Parameter{
		{
			Name:        "src",
			Type:        cty.String,
			Description: "The source code of the HCL native syntax configuration to parse and evaluate.",
		},
		{
			Name:        "vars",
			Type:        cty.DynamicPseudoType,
			Description: "An object describing the variables to include in the evaluation scope.",
		},
	},
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		src := args[0].AsString()
		varsVal := args[1]
		if !varsVal.Type().IsObjectType() {
			return cty.DynamicVal, function.NewArgErrorf(1, "must be an object whose attributes represent the variables to include in the HCL evaluation scope")
		}
		vars := varsVal.AsValueMap()

		file, diags := hclsyntax.ParseConfig([]byte(src), "<src>", hcl.InitialPos)
		if diags.HasErrors() {
			return cty.DynamicVal, function.NewArgErrorf(0, "invalid syntax: %s", diags.Error())
		}

		attrs, diags := file.Body.JustAttributes()
		if diags.HasErrors() {
			return cty.DynamicVal, function.NewArgErrorf(0, "invalid config: %s", diags.Error())
		}
		names := make([]string, 0, len(attrs))
		for name := range attrs {
			names = append(names, name)
		}
		sort.Strings(names)

		retVals := make(map[string]cty.Value, len(attrs))
		for _, name := range names {
			attr := attrs[name]
			v, diags := attr.Expr.Value(&hcl.EvalContext{
				Variables: vars,
				Functions: evalScopeFuncs,
			})
			if diags.HasErrors() {
				return cty.DynamicVal, function.NewArgErrorf(0, "evaluation failed for %q: %s", name, diags.Error())
			}
			retVals[name] = v
		}

		return cty.ObjectVal(retVals), nil
	},
}
