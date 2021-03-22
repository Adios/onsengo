// Package expression wraps a `goja` vm to run a given JavaScript expression in a string.
package expression

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

type Expression interface {
	// Run the given JavaScript code which produces a *value*, i.e. expressions.
	// Returns a string of the value's JSON representation and any JS error encountered,
	// such as exceptions, syntax errors and expressions that were evaluated to `undefined`.
	Stringify(string) (string, error)
}

type expression struct {
	vm *goja.Runtime
}

func New() Expression {
	return &expression{}
}

func (e *expression) Stringify(js string) (string, error) {
	torun := fmt.Sprintf("JSON.stringify(%s)", js)

	res, err := e.getVm().RunString(torun)
	if err != nil {
		return "", err
	}

	out := res.Export()
	if out == nil {
		err = errors.New(
			"We got nothing after running. Possibly the js returned an undefined.",
		)
		return "", err
	}

	return out.(string), nil
}

func (e *expression) getVm() *goja.Runtime {
	if e.vm == nil {
		e.vm = goja.New()
	}
	return e.vm
}
