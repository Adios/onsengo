// Package html provides utilities to handle web page https://onsen.ag/.
package html

import (
	"regexp"
)

// Parse the html content and return a pointer to the capture of first appeared NUXT data object:
//   <script>window.__NUXT__=([^<]+);</script>
// Returns nil if not found.
func FindNuxtExpression(html string) *string {
	re := regexp.MustCompile("<script>window.__NUXT__=([^<]+);</script>")

	m := re.FindStringSubmatch(html)
	if m == nil {
		return nil
	}
	return &m[1]
}
