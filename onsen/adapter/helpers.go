package adapter

import (
	"regexp"
	"strconv"
	"time"
)

// Given an incomplete date string in the format "MM/DD", and with a referenced date.
// The function finds a nearest year such that "YYYY/MM/DD" doesn't go over the referenced date.
func GuessTime(guess string, reference time.Time) time.Time {
	re := regexp.MustCompile("^([0-9]{1,2})/([0-9]{1,2})$")
	m := re.FindStringSubmatch(guess)

	if m == nil {
		return time.Time{}
	}

	guessMonth, err := strconv.Atoi(m[1])
	if err != nil {
		panic(err)
	}
	guessDay, err := strconv.Atoi(m[2])
	if err != nil {
		panic(err)
	}

	attemptTime := time.Date(
		reference.Year(),
		time.Month(guessMonth),
		guessDay, 0, 0, 0, 0,
		reference.Location(),
	)

	if attemptTime.After(reference) {
		return attemptTime.AddDate(-1, 0, 0)
	} else {
		return attemptTime
	}
}

// Based on the GuessTime() function, here we set a timezone of UTC+9 and use time.Now() as a referenced date.
func GuessJstTimeWithNow(guess string) time.Time {
	loc := time.FixedZone("UTC+9", 9*60*60)
	now := time.Now().In(loc)

	return GuessTime(guess, now)
}
