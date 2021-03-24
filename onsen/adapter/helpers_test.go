package adapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGuessTime(t *testing.T) {
	ref := time.Date(2021, time.Month(3), 24, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, time.Time{}, GuessTime("2020/202/30", ref))

	expected := time.Date(2020, time.Month(3), 25, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, GuessTime("3/25", ref))

	assert.Equal(t, ref.AddDate(0, 0, -1), GuessTime("3/23", ref))
	assert.Equal(t, ref, GuessTime("3/24", ref))
}
