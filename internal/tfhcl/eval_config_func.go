package tfhcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/apparentlymart/terraform-provider-hcl/internal/speclang"
)

var evalconfigFunc = &function.Spec{
	Description: "Evaluates HCL configuration source code in \"Just Attributes\" mode, returning an object representing the attributes.",
	Params: []function.Parameter{
		{
			Name:        "src",
			Type:        cty.String,
			Description: "The source code of the HCL native syntax configuration to parse and evaluate.",
		},
		{
			Name:        "specfile",
			Type:        cty.String,
			Description: "Path to the file containing the hcldec specification describing how to interpret the configuration.",
		},
		{
			Name:        "vars",
			Type:        cty.DynamicPseudoType,
			Description: "An object describing any additional variables to include in the evaluation scope.",
		},
	},
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		src := args[0].AsString()
		specFilename := args[1].AsString()
		varsVal := args[2]
		if !varsVal.Type().IsObjectType() {
			return cty.DynamicVal, function.NewArgErrorf(2, "must be an object whose attributes represent the variables to include in the HCL evaluation scope")
		}
		vars := varsVal.AsValueMap()

		specFile, diags := speclang.LoadSpecFile(specFilename)
		if diags.HasErrors() {
			return cty.DynamicVal, function.NewArgErrorf(1, "invalid spec file: %s", diags.Error())
		}

		file, diags := hclsyntax.ParseConfig([]byte(src), "<src>", hcl.InitialPos)
		if diags.HasErrors() {
			return cty.DynamicVal, function.NewArgErrorf(0, "invalid syntax: %s", diags.Error())
		}

		baseCtx := &hcl.EvalContext{
			Variables: specFile.Variables,
			Functions: specFile.Functions,
		}
		evalCtx := baseCtx
		if len(vars) != 0 {
			evalCtx = baseCtx.NewChild()
			evalCtx.Variables = vars
		}

		ret, diags := hcldec.Decode(file.Body, specFile.RootSpec, evalCtx)
		if diags.HasErrors() {
			return cty.DynamicVal, function.NewArgErrorf(0, "evaluation failed: %s", diags.Error())
		}

		return ret, nil
	},
}
