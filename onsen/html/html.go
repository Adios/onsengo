// Package html provides an abstract layer for the web page https://onsen.ag/.
package html

import (
	"regexp"
)

type Html struct {
	html  string
	nuxt  *string
	state int
}

const (
	hasntParsed = iota
	hasParsed
)

// Takes a string containing html content.
func New(html string) Html {
	return Html{
		html: html,
	}
}

// Parse the html content and return a pointer to the capture of first appeared NUXT data object:
//   <script>window.__NUXT__=([^<]+;)</script>
// Returns nil if not found. Parsed result will be cached in the (h *Html).
func (h *Html) GetNuxtExpression() *string {
	if h.state == hasntParsed {
		h.state = hasParsed

		re := regexp.MustCompile("<script>window.__NUXT__=([^<]+;)</script>")

		m := re.FindStringSubmatch(h.html)
		if m != nil {
			h.nuxt = &m[1]
		}
	}

	return h.nuxt
}
