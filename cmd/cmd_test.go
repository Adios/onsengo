package cmd

import (
	"compress/bzip2"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

type mockEpisode struct {
	manifest string
	tm       time.Time
}

func (e mockEpisode) Id() int {
	return 1227
}

func (e mockEpisode) Manifest() (string, bool) {
	if e.manifest == "" {
		return "", false
	}
	return e.manifest, true
}

func (e mockEpisode) JstUpdatedAt() (time.Time, bool) {
	if e.tm == (time.Time{}) {
		return time.Time{}, false
	}
	return e.tm, true
}

func TestFilter(t *testing.T) {
	var (
		assert = assert.New(t)
		pt     = func(str string) time.Time {
			out, _ := time.Parse("2006-01-02", str)
			return out
		}
		empty    = mockEpisode{}
		withMani = mockEpisode{manifest: "a"}
		withTime = mockEpisode{tm: pt("1989-12-27")}
		normal   = mockEpisode{manifest: "b", tm: pt("2021-06-04")}
	)

	{
		f := NewFilter()
		f.Push(empty)
		f.Push(withMani)
		f.Push(withTime)
		f.Push(normal)
		assert.Equal([]string{"a", "b"}, f.Out())
	}
	{
		out := strings.Builder{}
		log.SetOutput(&out)

		f := NewFilter(FilterUpdatedAfter(pt("2021-06-04")))
		f.Push(empty)
		assert.Contains(out.String(), "doesn't have update time")
		out.Reset()
		f.Push(withMani)
		assert.Contains(out.String(), "doesn't have update time")
		f.Push(withTime)
		f.Push(normal)
		assert.Equal([]string{"b"}, f.Out())
	}
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

	execute(func(out b, err b) {
		assert.NoError(Execute())
		assert.Equal(strings.Repeat("HAS_BEEN_SCREENED\n", 121), out.String())
	}, "lsm", "--backend", server.URL)

	execute(func(out b, err b) {
		assert.NoError(Execute())
		assert.Equal(strings.Repeat("HAS_BEEN_SCREENED\n", 2), out.String())
		assert.Equal("nosuchradio: not found\n", err.String())
	}, "lsm", "fujita", "gurepap", "nosuchradio", "--backend", server.URL)

	execute(func(out b, err b) {
		assert.NoError(Execute())
		assert.Equal("HAS_BEEN_SCREENED\n", out.String())
		assert.Equal("fujita/3560: empty manifest, may be inaccessible\nfujita/9999: not found\n", err.String())
	}, "lsm", "fujita/3598", "fujita/3560", "fujita/9999", "--backend", server.URL)

	execute(func(out b, err b) {
		assert.NoError(Execute())
		assert.Equal("HAS_BEEN_SCREENED\n", out.String())
		assert.Equal("fujita/3560: empty manifest, may be inaccessible\n", err.String())
	}, "lsm", "toshitai", "fujita/3598", "fujita/3560", "--after", "2021-03-16", "--backend", server.URL)

	execute(func(out b, err b) {
		var (
			f, _    = os.Open("testdata/expected_dump.txt.bz2")
			r       = bzip2.NewReader(f)
			data, _ = io.ReadAll(r)
		)

		assert.NoError(Execute())
		assert.Equal(string(data), out.String())
	}, "dump", "--backend", server.URL)
}

func server(t *testing.T) http.Handler {
	f, _ := os.ReadFile("testdata/fixture_nologin_screened.html")

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, ua, req.Header.Get("User-Agent"))
		w.Write(f)
	})
}
