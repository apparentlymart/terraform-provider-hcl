package tfhcl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestEvalconfigattrs(t *testing.T) {
	tests := []struct {
		Expr    string
		Vars    map[string]cty.Value
		Want    cty.Value
		WantErr string
	}{
		{
			Expr: `a = 1`,
			Want: cty.ObjectVal(map[string]cty.Value{
				"a": cty.NumberIntVal(1),
			}),
		},
		{
			Expr: ``,
			Want: cty.EmptyObjectVal,
		},
		{
			Expr: `a = unk + 2`,
			Vars: map[string]cty.Value{
				"unk": cty.UnknownVal(cty.String),
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"a": cty.UnknownVal(cty.Number).RefineNotNull(),
			}),
		},
		{
			Expr:    `invalid syntax`,
			WantErr: `invalid syntax: <src>:1,15-15: Invalid block definition; Either a quoted string block label or an opening brace ("{") is expected here.`,
		},
		{
			Expr:    `a = {} + 1`,
			WantErr: `evaluation failed for "a": <src>:1,5-7: Invalid operand; Unsuitable value for left operand: number required.`,
		},
	}

	p := NewProvider()
	f := p.CallStub("evalconfigattrs")
	for _, test := range tests {
		t.Run(test.Expr, func(t *testing.T) {
			got, gotErr := f(cty.StringVal(test.Expr), cty.ObjectVal(test.Vars))

			if test.WantErr != "" {
				if gotErr == nil {
					t.Fatalf("unexpected success\nwant error: %s", test.WantErr)
				}
				if got, want := gotErr.Error(), test.WantErr; got != want {
					t.Fatalf("wrong error\ngot:  %s\nwant: %s", got, want)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("unexpected error\ngot error: %s\nwant: %#v", gotErr.Error(), test.Want)
			}
			if diff := cmp.Diff(test.Want, got, ctydebug.CmpOptions); diff != "" {
				t.Errorf("wrong result\nexpr: %s\n%s", test.Expr, diff)
			}
		})
	}

}
