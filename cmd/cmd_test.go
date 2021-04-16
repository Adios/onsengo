package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/adios/onsengo/onsen"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set a fixed date instead of time.Now() in JstUpdatedAt()
	onsen.SetRefDate("2021-04-14")
	// Shortcut both onsen & cobra's output/stderr
	root.out, root.err = new(strings.Builder), new(strings.Builder)
	root.cmd.SetOut(root.out)
	root.cmd.SetErr(root.err)

	os.Exit(m.Run())
}

func TestUnique(t *testing.T) {
	type s = []string

	tests := []struct {
		in       s
		expected s
	}{
		{s{}, s{}},
		{s{"a"}, s{"a"}},
		{s{"a", "a"}, s{"a"}},
		{s{"a", "b"}, s{"a", "b"}},
		{s{"a", "b", "a", "b", "c"}, s{"a", "b", "c"}},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, unique(test.in))
	}
}

func Test(t *testing.T) {
	type b = *strings.Builder

	var (
		assert = assert.New(t)
		server = httptest.NewServer(server(t))

		execute = func(fn func(b, b), input ...string) {
			out := root.out.(b)
			err := root.err.(b)

			root.cmd.SetArgs(input)
			fn(out, err)

			out.Reset()
			err.Reset()
		}
	)
	defer server.Close()

	execute(func(out b, err b) {
		assert.EqualError(Execute(), "Create: NUXT pattern not matched")
	}, "ls", "--backend", "file:///")

	execute(func(out b, err b) {
		f, _ := os.ReadFile("testdata/expected_ls.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
	}, "ls", "--backend", server.URL)

	execute(func(out b, err b) {
		f, _ := os.ReadFile("testdata/expected_ls_recursive.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
	}, "ls", "-r", "--backend", server.URL)

	execute(func(out b, err b) {
		f, _ := os.ReadFile("testdata/expected_ls_single.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
		assert.Equal("nosuchradio: not found\n", err.String())
	}, "ls", "fujita", "gurepap", "nosuchradio", "--backend", server.URL)

	execute(func(out b, err b) {
		f, _ := os.ReadFile("testdata/expected_ls_single.txt")
		assert.NoError(Execute())
		assert.Equal(string(f), out.String())
		assert.Equal("nosuchradio: not found\n", err.String())
	}, "ls", "fujita", "gurepap", "nosuchradio", "fujita", "gurepap", "--backend", server.URL)
}

func server(t *testing.T) http.Handler {
	f, _ := os.ReadFile("testdata/fixture_nologin_screened.html")

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, ua, req.Header.Get("User-Agent"))
		w.Write(f)
	})
}
