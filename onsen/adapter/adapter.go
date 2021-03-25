// Package adapter defines the interface to interact with data adapters.
package adapter

import "time"

type Adapter interface {
	RadioShows() []RadioShow
	// Returns nil if there is no login associated.
	User() User
}

// A signed-in user.
type User interface {
	Email() string
	// A string of digits.
	UserId() string
	// A slice of uint returned from Person.PersonId()
	FollowingPeople() []uint
	// A slice of uint returned RadioShow.RadioShowId()
	FollowingRadioShows() []uint
	// A slice of uint returned Episode.EpisodeId()
	PlayingEpisodes() []uint
}

// A radio series.
type RadioShow interface {
	RadioShowId() uint
	Name() string
	Title() string
	HasUpdates() bool

	// Returns a best-effor time that is guessed based on time.Now().
	// Since there is no YYYY recorded in onsen's raw data. (MM/DD only)
	// An empty time.Time{} means there is an invalid date pattern.
	GuessedUpdatedAt() time.Time

	Hosts() []Person

	// Returns a slice of Episode instances which may either be an AudioEpisode or a VideoEpisode.
	Episodes() []Episode
}

// Personality.
type Person interface {
	PersonId() uint
	Name() string
}

//An episode of a radio series.
type Episode interface {
	EpisodeId() uint
	RadioShowId() uint
	Title() string

	// Returns the url to the episode's poster image.
	Poster() string

	// Returns the url to the episode's manifest (m3u8).
	// An empty string ("") means the resource is not accessible from current user identity.
	Manifest() string

	// Returns a best-effor time that is guessed based on time.Now().
	// Since there is no YYYY recorded in onsen's raw data. (MM/DD only)
	// An empty time.Time{} means there is an invalid date pattern.
	GuessedPublishedAt() time.Time

	// Always returns a slice, an empty slice means there are no guests.
	//  [ "name1", "name2" ]
	Guests() []string

	IsBonus() bool
	IsSticky() bool
	IsLatest() bool
	RequiresPremium() bool
}

// Audio is attached to the Episodes that deliver sound-only contents.
// No methods for now
type Audio interface {
	Audio()
}

// Video is attached to the Episodes that deliver video contents.
// No methods for now
type Video interface {
	Video()
}
