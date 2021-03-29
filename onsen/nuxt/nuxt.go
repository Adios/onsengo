// Package nuxt decodes a JSON.stringified Nuxt object into a struct.
package nuxt

import (
	"encoding/json"
	"io"
)

// Represents the root of a Nuxt JSON object. Decodes only the fields we want.
type Root struct {
	Error     interface{} `json:"error"`
	State     State       `json:"state"`
	RoutePath string      `json:"routePath"`
}

// Represents the root.State of a Nuxt JSON object. Decodes only the fields we want.
// If it is an anonymous session, Signin will be nil.
type State struct {
	Signin   *Signin `json:"sign_in"`
	Programs struct {
		Programs struct {
			All []Program `json:"all"`
		} `json:"programs"`
	} `json:"programs"`
}

// Represents the root.state.sign_in of a Nuxt JSON object. Decodes only the fields we want.
type Signin struct {
	Email                string `json:"email"`
	Id                   string `json:"id"`
	FavoritePerformerIds []uint `json:"favorite_performer_ids"`
	FavoriteProgramIds   []uint `json:"favorite_program_ids"`
	PlaylistedContentIds []uint `json:"playlisted_content_ids"`
}

// Represents the root.state.programs.programs.all[] of a Nuxt JSON object. Decodes only the fields we want.
// If a radio series is got announced and it has no contents, as well as some special programs, they will have nil Updated.
type Program struct {
	Id            uint        `json:"id"`
	DirectoryName string      `json:"directory_name"`
	Title         string      `json:"title"`
	New           bool        `json:"new"`
	Updated       *string     `json:"updated"`
	Performers    []Performer `json:"performers"`
	Contents      []Content   `json:"contents"`
}

// Represents the root.state.programs.programs.all[].performers of a Nuxt JSON object. Decodes all fields.
type Performer struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// Represents the root.state.programs.programs.all[].Contents[] of a Nuxt JSON object. Decodes only the fields we want.
// If the current user identity (or anonymous) has no permissions to play the content, StreamingUrl will be nil.
type Content struct {
	Id             uint     `json:"id"`
	Title          string   `json:"title"`
	Bonus          bool     `json:"bonus"`
	Sticky         bool     `json:"sticky"`
	Latest         bool     `json:"latest"`
	MediaType      string   `json:"media_type"`
	Premium        bool     `json:"premium"`
	ProgramId      uint     `json:"program_id"`
	DeliveryDate   string   `json:"delivery_date"`
	Movie          bool     `json:"movie"`
	PosterImageUrl string   `json:"poster_image_url"`
	StreamingUrl   *string  `json:"streaming_url"`
	Guests         []string `json:"guests"`
}

func FromReader(r io.Reader) (*Root, error) {
	var data Root

	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
