package extractor

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errMsg1 = errors.New("NUXT data not found")
)

func TestNoNuxtString(t *testing.T) {
	e := New("hello world")
	_, err := e.Extract()

	assert.Equal(t, errMsg1, err)
}

func TestInvalidNuxtString(t *testing.T) {
	e := New("<script>window.__NUXT__=...;</script>")
	_, err := e.Extract()

	assert.NotNil(t, err)
}

func TestOk(t *testing.T) {
	data, _ := os.ReadFile("test_fixture_onsen_ag.html")
	e := New(string(data))

	str, err := e.Extract()

	assert.Nil(t, err)
	assert.Equal(t, len(str), 1830323)
}

func Example() {
	// `curl -L https://onsen.ag > test_fixture_onsen_ag.html`, contains a obfuscated NUXT data string
	data, _ := os.ReadFile("test_fixture_onsen_ag.html")
	e := New(string(data))

	str, err := e.Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(str[0:42])
	// Output: {"layout":"default","data":[{"category":5}
}
