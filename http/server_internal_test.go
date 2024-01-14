package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/storvik/goshrt"
)

func newMockServer(tb testing.TB, m *mockStore) *Server {
	tb.Helper() // This function is a test-helper, not a test

	l := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	s := NewServer(l, "test")
	s.InfoLog = l
	s.ShrtStore = m
	s.Auth = &mockAuthorizer{}

	return s
}

// TODO: Add more tests with invalid requests id etc.
func TestShrtGetHandler(t *testing.T) {
	var tests = []struct {
		shrt               *goshrt.Shrt
		url                string
		expectedStatusCode int
	}{
		{
			shrt:               &goshrt.Shrt{ID: 1},
			url:                "/api/shrt/2",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			shrt: &goshrt.Shrt{
				ID:     1,
				Domain: "domain",
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/api/shrt/1",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		// Create mock store and populate it with shrt.
		m := &mockStore{shrt: tt.shrt}

		s := newMockServer(t, m)
		req := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)

		w := executeRequest(req, s)

		res := w.Result()
		resShrt := new(goshrt.Shrt)
		decoder := json.NewDecoder(res.Body)

		err := decoder.Decode(&resShrt)
		if err != nil {
			t.Errorf("error decoding response, %v", err)
		}

		checkResponseCode(t, tt.expectedStatusCode, res.StatusCode)

		// Check if received shrt is correct, but only if http status is ok.
		if res.StatusCode == http.StatusOK {
			if !reflect.DeepEqual(tt.shrt, resShrt) {
				t.Errorf("shorts not equal, expected %v, got %v", tt.shrt, resShrt)
			}
		}

		res.Body.Close()
	}
}

func TestShrtDeleteHandler(t *testing.T) {
	var tests = []struct {
		shrt               *goshrt.Shrt
		url                string
		expectedStatusCode int
	}{
		{
			shrt:               &goshrt.Shrt{ID: 1},
			url:                "/api/shrt/2",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			shrt: &goshrt.Shrt{
				ID:     1,
				Domain: "domain",
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/api/shrt/1",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		// Create mock store and populate it with shrt.
		m := &mockStore{shrt: tt.shrt}

		s := newMockServer(t, m)
		req := httptest.NewRequest(http.MethodDelete, tt.url, http.NoBody)

		w := executeRequest(req, s)

		res := w.Result()
		resShrt := new(goshrt.Shrt)
		decoder := json.NewDecoder(res.Body)

		err := decoder.Decode(&resShrt)
		if err != nil {
			t.Errorf("error decoding response, %v", err)
		}

		checkResponseCode(t, tt.expectedStatusCode, res.StatusCode)

		// Check if received shrt is correct, but only if http status is ok.
		if res.StatusCode == http.StatusOK {
			if !reflect.DeepEqual(tt.shrt, resShrt) {
				t.Errorf("shorts not equal, expected %v, got %v", tt.shrt, resShrt)
			}
		}

		res.Body.Close()
	}
}
