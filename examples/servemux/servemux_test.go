package servemux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "users")
	})
	mux.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "user")
	})
	return mux
}

func TestServeMuxRoutes(t *testing.T) {
	t.Parallel()
	mux := newMux()

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	cases := []struct {
		path string
		want string
	}{
		{"/api/users/", "users"},
		{"/api/users/123", "user"},
	}
	for _, c := range cases {
		resp, err := ts.Client().Get(ts.URL + c.path)
		if err != nil {
			t.Fatalf("GET %s: %v", c.path, err)
		}
		defer resp.Body.Close()
		buf := make([]byte, 16)
		n, _ := resp.Body.Read(buf)
		got := string(buf[:n])
		if got != c.want {
			t.Fatalf("%s: got=%q want=%q", c.path, got, c.want)
		}
	}
}
