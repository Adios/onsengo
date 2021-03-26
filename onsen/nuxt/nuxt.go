// Package nuxt represents the structure of a Nuxt object in https://onsen.ag/.
package nuxt

import (
	"encoding/json"
)

// Root structure of the Nuxt object. Only the fields we may use were defined here.
type Nuxt struct {
	Error     interface{} `json:"error"`
	State     State       `json:"state"`
	RoutePath string      `json:"routePath"`
}

// Only the fields we may use were defined here.
type State struct {
	Signin   *Signin `json:"sign_in"`
	Programs struct {
		Programs struct {
			All []Program `json:"all"`
		} `json:"programs"`
	} `json:"programs"`
}

// Personal data. Only the fields we may use were defined here.
type Signin struct {
	Email                string        `json:"email"`
	Id                   string        `json:"id"`
	FavoritePerformerIds []PerformerId `json:"favorite_performer_ids"`
	FavoriteProgramIds   []ProgramId   `json:"favorite_program_ids"`
	PlaylistedContentIds []ContentId   `json:"playlisted_content_ids"`
}

// A radio program. Only the fields we may use were defined here.
type Program struct {
	Id            ProgramId   `json:"id"`
	DirectoryName string      `json:"directory_name"`
	Title         string      `json:"title"`
	New           bool        `json:"new"`
	Updated       *string     `json:"updated"`
	Performers    []Performer `json:"performers"`
	Contents      []Content   `json:"contents"`
}

// A performer.
type Performer struct {
	Id   PerformerId `json:"id"`
	Name string      `json:"name"`
}

// The episode of a radio program. Only the fields we may use were defined here.
type Content struct {
	Id             ContentId `json:"id"`
	Title          string    `json:"title"`
	Bonus          bool      `json:"bonus"`
	Sticky         bool      `json:"sticky"`
	Latest         bool      `json:"latest"`
	MediaType      string    `json:"media_type"`
	Premium        bool      `json:"premium"`
	ProgramId      ProgramId `json:"program_id"`
	DeliveryDate   string    `json:"delivery_date"`
	Movie          bool      `json:"movie"`
	PosterImageUrl string    `json:"poster_image_url"`
	StreamingUrl   *string   `json:"streaming_url"`
	Guests         []string  `json:"guests"`
}

type PerformerId uint
type ProgramId uint
type ContentId uint

// Takes a JSON string, returns an unmarshaled data and any error encountered.
func Parse(str string) (*Nuxt, error) {
	var n Nuxt

	err := json.Unmarshal([]byte(str), &n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}
