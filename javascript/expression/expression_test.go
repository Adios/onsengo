package expression

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestUndefined(t *testing.T) {
	e := New()
	res, err := e.Stringify("")

	assert.Equal(t, res, "")
	assert.Equal(t, err, errors.New("We got nothing after running. Possibly the js returned an undefined."))
}

func TestSyntaxError(t *testing.T) {
	e := New()
	res, err := e.Stringify(";")

	assert.Equal(t, res, "")

	_, ok := err.(*goja.Exception)
	assert.Equal(t, ok, true)
}

func Example() {
	e := New()
	res, err := e.Stringify("(function(x) {return {hello: x}})('world')")

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	// Output: {"hello":"world"}
}
