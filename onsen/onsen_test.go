package onsen

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdConversion(t *testing.T) {
	assert.Equal(t, "0", PersonId(0).String())
	assert.Equal(t, "55", PersonId(55).String())
	assert.Equal(t, "555", RadioShowId(555).String())
	assert.Equal(t, "5555", EpisodeId(5555).String())
}

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
