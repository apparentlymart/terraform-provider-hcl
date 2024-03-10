package tfhcl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestEvalexpr(t *testing.T) {
	tests := []struct {
		Expr    string
		Vars    map[string]cty.Value
		Want    cty.Value
		WantErr string
	}{
		{
			Expr: `1`,
			Want: cty.NumberIntVal(1),
		},
		{
			Expr: `"hello"`,
			Want: cty.StringVal("hello"),
		},
		{
			Expr: `upper("hello")`,
			Want: cty.StringVal("HELLO"),
		},
		{
			Expr: `foo`,
			Vars: map[string]cty.Value{
				"foo": cty.StringVal("hello"),
			},
			Want: cty.StringVal("hello"),
		},
		{
			Expr: `unk + 2`,
			Vars: map[string]cty.Value{
				"unk": cty.UnknownVal(cty.String),
			},
			Want: cty.UnknownVal(cty.Number).RefineNotNull(),
		},
		{
			Expr:    `invalid syntax`,
			WantErr: `invalid syntax: <src>:1,9-15: Extra characters after expression; An expression was successfully parsed, but extra characters were found after it.`,
		},
		{
			Expr:    `{} + 1`,
			WantErr: `evaluation failed: <src>:1,1-3: Invalid operand; Unsuitable value for left operand: number required.`,
		},
	}

	p := NewProvider()
	f := p.CallStub("evalexpr")
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

func TestEvaltemplate(t *testing.T) {
	tests := []struct {
		Template string
		Vars     map[string]cty.Value
		Want     cty.Value
		WantErr  string
	}{
		{
			Template: `${1}`,
			Want:     cty.StringVal("1"),
		},
		{
			Template: `${"hello"}`,
			Want:     cty.StringVal("hello"),
		},
		{
			Template: `${upper("hello")}`,
			Want:     cty.StringVal("HELLO"),
		},
		{
			Template: `${foo}`,
			Vars: map[string]cty.Value{
				"foo": cty.StringVal("hello"),
			},
			Want: cty.StringVal("hello"),
		},
		{
			Template: `${unk + 2}`,
			Vars: map[string]cty.Value{
				"unk": cty.UnknownVal(cty.String),
			},
			Want: cty.UnknownVal(cty.String).RefineNotNull(),
		},
		{
			Template: `foo-${unk}`,
			Vars: map[string]cty.Value{
				"unk": cty.UnknownVal(cty.String),
			},
			Want: cty.UnknownVal(cty.String).Refine().
				NotNull().
				StringPrefixFull("foo-").
				NewValue(),
		},
		{
			Template: `${invalid syntax}`,
			WantErr: `invalid syntax: <src>:1,11-17: Extra characters after interpolation expression; Expected a closing brace to end the interpolation expression, but found extra characters.

This can happen when you include interpolation syntax for another language, such as shell scripting, but forget to escape the interpolation start token. If this is an embedded sequence for another language, escape it by starting with "$${" instead of just "${".`,
		},
		{
			Template: `${[] + 1}`,
			WantErr:  `evaluation failed: <src>:1,3-5: Invalid operand; Unsuitable value for left operand: number required.`,
		},
		{
			Template: `${[]}`,
			WantErr:  `invalid result type: string required`,
		},
	}

	p := NewProvider()
	f := p.CallStub("evaltemplate")
	for _, test := range tests {
		t.Run(test.Template, func(t *testing.T) {
			got, gotErr := f(cty.StringVal(test.Template), cty.ObjectVal(test.Vars))

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
				t.Errorf("wrong result\ntemplate: %s\n%s", test.Template, diff)
			}
		})
	}
}
