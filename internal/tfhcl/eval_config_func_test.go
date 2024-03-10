package tfhcl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestEvalconfig(t *testing.T) {
	tests := map[string]struct {
		Config   string
		SpecFile string
		Vars     map[string]cty.Value
		Want     cty.Value
		WantErr  string
	}{
		"empty": {
			Config:   ``,
			SpecFile: "testdata/simplespec.hcldec",
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.NullVal(cty.String),
			}),
		},
		"name constant": {
			Config:   `name = "Jackson"`,
			SpecFile: "testdata/simplespec.hcldec",
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jackson"),
			}),
		},
		"name variable": {
			Config:   `name = foo`,
			SpecFile: "testdata/simplespec.hcldec",
			Vars: map[string]cty.Value{
				"foo": cty.StringVal("Jackson"),
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jackson"),
			}),
		},
		"name unknown": {
			Config:   `name = foo`,
			SpecFile: "testdata/simplespec.hcldec",
			Vars: map[string]cty.Value{
				"foo": cty.DynamicVal,
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.UnknownVal(cty.String),
			}),
		},
		"name unknown but refined": {
			Config:   `name = foo`,
			SpecFile: "testdata/simplespec.hcldec",
			Vars: map[string]cty.Value{
				"foo": cty.UnknownVal(cty.String).Refine().
					NotNull().
					StringPrefix("boop-").
					NewValue(),
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.UnknownVal(cty.String).Refine().
					NotNull().
					StringPrefix("boop-").
					NewValue(),
			}),
		},
		"unexpected argument": {
			Config:   `nome = "Foo"`,
			SpecFile: "testdata/simplespec.hcldec",
			WantErr:  `evaluation failed: <src>:1,1-5: Unsupported argument; An argument named "nome" is not expected here. Did you mean "name"?`,
		},
		"invalid syntax in src": {
			Config:   `invalid syntax`,
			SpecFile: "testdata/simplespec.hcldec",
			WantErr:  `invalid syntax: <src>:1,15-15: Invalid block definition; Either a quoted string block label or an opening brace ("{") is expected here.`,
		},
		"invalid syntax in spec": {
			Config:   ``,
			SpecFile: "testdata/badsyntax.hcldec",
			WantErr:  `invalid spec file: testdata/badsyntax.hcldec:1,15-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.`,
		},
	}

	p := NewProvider()
	f := p.CallStub("evalconfig")
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := f(cty.StringVal(test.Config), cty.StringVal(test.SpecFile), cty.ObjectVal(test.Vars))

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
				t.Errorf("wrong result\n\n%s", diff)
			}
		})
	}
}
