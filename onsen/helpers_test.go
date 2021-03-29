package onsen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyString(t *testing.T) {
	expr, ok := FindNuxtExpression("")

	assert.Empty(t, expr)
	assert.False(t, ok)
}

func TextNoNuxtPattern(t *testing.T) {
	html := "...<script>window.__NUXT__=one</script><script>window.__NUXT__=two</script>..."
	expr, ok := FindNuxtExpression(html)

	assert.Nil(t, expr)
	assert.False(t, ok)
}

func TextNuxtPattern(t *testing.T) {
	html := "...<script>window.__NUXT__=one</script><script>window.__NUXT__=two;</script>..."
	expr, ok := FindNuxtExpression(html)

	assert.Equal(t, "two", string(expr))
	assert.True(t, ok)
}

func Example() {
	html := "...<script>window.__NUXT__=one;</script><script>window.__NUXT__=two;</script>..."
	expr, ok := FindNuxtExpression(html)
	if !ok {
		panic("not found")
	}

	fmt.Println(expr)
	// Output: one
}
