package tfhcl

import (
	"github.com/apparentlymart/go-tf-func-provider/tffunc"
)

func NewProvider() *tffunc.Provider {
	p := tffunc.NewProvider()
	p.AddFunction("evalconfig", evalconfigFunc)
	p.AddFunction("evalconfigattrs", evalconfigattrsFunc)
	p.AddFunction("evalexpr", evalexprFunc)
	p.AddFunction("evaltemplate", evaltemplateFunc)
	return p
}
