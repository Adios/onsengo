package expression

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestUndefined(t *testing.T) {
	e := From("")
	res, err := e.Stringify()

	assert.Equal(t, "", res)
	assert.Equal(t, fmt.Errorf("Got nothing after running. Possibly the js returned an undefined.\n"), err)
}

func TestSyntaxError(t *testing.T) {
	e := From(";")
	res, err := e.Stringify()

	assert.Equal(t, "", res)

	_, ok := err.(*goja.Exception)
	assert.True(t, ok)
}

func Example() {
	e := From("(function(x) {return {hello: x}})('world')")

	res, err := e.Stringify()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	// Output: {"hello":"world"}
}
