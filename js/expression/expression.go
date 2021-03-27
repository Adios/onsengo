// Package expression is a wrapper of github.com/dop251/goja to run Javascript code (expressions) for deobfuscation.
package expression

import (
	"fmt"

	"github.com/dop251/goja"
)

type Expression struct {
	js string
	vm *goja.Runtime
}

// Run the given JavaScript code which can produces a *value*, i.e. expressions.
// Returns a string of the value's JSON representation and any JS error encountered.
// Note that "undefined" is also considered as an error.
func (e *Expression) Stringify() (json string, err error) {
	torun := fmt.Sprintf("JSON.stringify(%s)", string(e.js))

	res, err := e.getVm().RunString(torun)
	if err != nil {
		return "", err
	}

	out := res.Export()
	if out == nil {
		return "", fmt.Errorf("Got nothing after running. Possibly the js returned an undefined.\n")
	}

	return out.(string), nil
}

func (e *Expression) getVm() *goja.Runtime {
	if e.vm == nil {
		e.vm = goja.New()
	}
	return e.vm
}

func From(js string) *Expression {
	return &Expression{
		js: js,
	}
}
