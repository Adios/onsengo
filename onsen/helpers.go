package onsen

import (
	"regexp"
)

// Returns a string to the capture of first appeared NUXT pattern:
//   <script>window.__NUXT__=([^<]+);</script>
func FindNuxtExpression(html string) (expr string, ok bool) {
	re := regexp.MustCompile("<script>window.__NUXT__=([^<]+);</script>")

	m := re.FindStringSubmatch(html)
	if m == nil {
		return "", false
	}
	return m[1], true
}
