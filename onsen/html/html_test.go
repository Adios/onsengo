package html

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyHtml(t *testing.T) {
	h := New("")
	res := h.GetNuxtExpression()

	assert.Nil(t, res)
}

func TestInvalidHtml(t *testing.T) {
	h := New("...<script>window.__NUXT__=one</script><script>window.__NUXT__=two</script>...")
	res := h.GetNuxtExpression()

	assert.Nil(t, res)
}

func TestMatchSecond(t *testing.T) {
	h := New("...<script>window.__NUXT__=one</script><script>window.__NUXT__=two;</script>...")
	res := h.GetNuxtExpression()

	assert.Equal(t, *res, "two;")
}

func Example() {
	h := New("...<script>window.__NUXT__=one;</script><script>window.__NUXT__=two;</script>...")
	res := h.GetNuxtExpression()

	if res == nil {
		panic(res)
	}

	fmt.Println(*res)
	// Output: one;
}
