package html

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyContent(t *testing.T) {
	expr := FindNuxtExpression("")

	assert.Nil(t, expr)
}

func TextNoNuxtPattern(t *testing.T) {
	expr := FindNuxtExpression("...<script>window.__NUXT__=one</script><script>window.__NUXT__=two</script>...")

	assert.Nil(t, expr)
}

func TextNuxtPattern(t *testing.T) {
	expr := FindNuxtExpression("...<script>window.__NUXT__=one</script><script>window.__NUXT__=two;</script>...")

	assert.Equal(t, "two", *expr)
}

func Example() {
	expr := FindNuxtExpression("...<script>window.__NUXT__=one;</script><script>window.__NUXT__=two;</script>...")
	if expr == nil {
		panic("not found")
	}

	fmt.Println(*expr)
	// Output: one
}
