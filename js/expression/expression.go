// Package expression wraps a `goja` vm to run a given JavaScript expression.
package expression

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

type Expression interface {
	// Run the contained JavaScript code which produces a *value*, i.e. expressions.
	// Returns a string of the value's JSON representation and any JS error encountered,
	// such as exceptions, syntax errors and expressions that were evaluated to `undefined`.
	Stringify() (string, error)
}

type expression struct {
	js string
	vm *goja.Runtime
}

// Takes a piece of JavaScript code to be evaluated.
func New(js string) Expression {
	return &expression{
		js: js,
	}
}

func (e *expression) Stringify() (string, error) {
	torun := fmt.Sprintf("JSON.stringify(%s)", e.js)

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
