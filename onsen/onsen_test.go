package onsen

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyString(t *testing.T) {
	onsen, err := Create("")

	assert.NotNil(t, err)
	assert.Nil(t, onsen)
}

func TestOK(t *testing.T) {
	content, err := os.ReadFile("fixture_onsen.html")
	if err != nil {
		panic(err)
	}

	onsen, err := Create(string(content))

	assert.Nil(t, err)

	assert.Equal(t, 128, len(onsen.RadioShows()))

	fmt.Printf("{%#v}", onsen.RadioShows())
}
