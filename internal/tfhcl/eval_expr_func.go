package tfhcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/apparentlymart/terraform-provider-hcl/internal/evalfuncs"
)

var evalexprFunc = parseAndEvalFunctionSpec(
	"Evaluates a given string as an HCL expression.",
	"The source code of the expression to evaluate.",
	hclsyntax.ParseExpression,
	cty.DynamicPseudoType,
)
var evaltemplateFunc = parseAndEvalFunctionSpec(
	"Evaluates a given string as an HCL template.",
	"The source code of the template to evaluate.",
	hclsyntax.ParseTemplate,
	cty.String,
)

func parseAndEvalFunctionSpec(desc, srcDesc string, parse func([]byte, string, hcl.Pos) (hclsyntax.Expression, hcl.Diagnostics), retType cty.Type) *function.Spec {
	ret := &function.Spec{
		Description: desc,
		Params: []function.Parameter{
			{
				Name:        "src",
				Type:        cty.String,
				Description: srcDesc,
			},
			{
				Name:        "vars",
				Type:        cty.DynamicPseudoType,
				Description: "An object describing the variables to include in the evaluation scope.",
			},
		},
		Type: function.StaticReturnType(retType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			src := args[0].AsString()
			varsVal := args[1]
			if !varsVal.Type().IsObjectType() {
				return cty.DynamicVal, function.NewArgErrorf(1, "must be an object whose attributes represent the variables to include in the HCL evaluation scope")
			}
			vars := varsVal.AsValueMap()

			expr, diags := parse([]byte(src), "<src>", hcl.InitialPos)
			if diags.HasErrors() {
				return cty.DynamicVal, function.NewArgErrorf(0, "invalid syntax: %s", diags.Error())
			}

			result, diags := expr.Value(&hcl.EvalContext{
				Variables: vars,
				Functions: evalfuncs.Functions(),
			})
			if diags.HasErrors() {
				return cty.DynamicVal, function.NewArgErrorf(0, "evaluation failed: %s", diags.Error())
			}

			result, err := convert.Convert(result, retType)
			if err != nil {
				return cty.DynamicVal, function.NewArgErrorf(0, "invalid result type: %s", err)
			}

			return result, nil
		},
	}
	if retType == cty.String {
		ret.RefineResult = func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		}
	}
	return ret
}
