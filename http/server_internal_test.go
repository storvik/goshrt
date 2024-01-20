package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

// TODO: Should probably add expire time to slugs.

func TestShrtHandler(t *testing.T) {
	var tests = []struct {
		shrt               *goshrt.Shrt
		url                string
		expectedStatusCode int
	}{
		{
			shrt: &goshrt.Shrt{
				ID:     1,
				Domain: "example.com", // example.com will be used as r.Host when using httptest.
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/slug",
			expectedStatusCode: http.StatusMovedPermanently,
		},
		{
			shrt: &goshrt.Shrt{
				ID:     1,
				Domain: "example.com", // example.com will be used as r.Host when using httptest.
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/notexistingslug",
			expectedStatusCode: http.StatusOK, // expecting status ok due to landingpage
		},
	}

	for _, tt := range tests {
		// Create mock store and populate it with shrt.
		m := &mockStore{shrt: tt.shrt}

		s := newMockServer(t, m)
		req := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)

		w := executeRequest(req, s)

		res := w.Result()

		checkResponseCode(t, tt.expectedStatusCode, res.StatusCode)

		res.Body.Close()
	}
}

func TestShrtCreateHandler(t *testing.T) {
	var tests = []struct {
		shrt               *goshrt.Shrt
		url                string
		expectedStatusCode int
	}{
		{
			shrt: &goshrt.Shrt{
				Domain: "error", // mockStore create returns error if domain is error
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/api/shrt",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			shrt: &goshrt.Shrt{
				Domain: "domain",
				Slug:   "slug",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/api/shrt",
			expectedStatusCode: http.StatusCreated,
		},
		{
			shrt:               &goshrt.Shrt{},
			url:                "/api/shrt",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		// Create mock store and populate it with shrt.
		m := &mockStore{shrt: tt.shrt}

		s := newMockServer(t, m)

		postBody, err := json.Marshal(tt.shrt)
		if err != nil {
			t.Errorf("error marshalling request, %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(postBody))

		w := executeRequest(req, s)

		res := w.Result()
		resShrt := new(goshrt.Shrt)
		decoder := json.NewDecoder(res.Body)

		err = decoder.Decode(&resShrt)
		if err != nil {
			t.Errorf("error decoding response, %v", err)
		}

		checkResponseCode(t, tt.expectedStatusCode, res.StatusCode)

		// Check if received shrt is correct, but only if http status is ok.
		if res.StatusCode == http.StatusCreated {
			compareShrts(t, resShrt, tt.shrt, false)
		}

		res.Body.Close()
	}
}

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
		{
			shrt: &goshrt.Shrt{
				ID:     2,
				Domain: "domaintest",
				Slug:   "slugtest",
				Dest:   "https://github.com/storvik/goshrt",
			},
			url:                "/api/shrt/domaintest/slugtest",
			expectedStatusCode: http.StatusOK,
		},
		{
			shrt:               &goshrt.Shrt{},
			url:                "/api/shrt/3",
			expectedStatusCode: http.StatusInternalServerError,
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
			compareShrts(t, resShrt, tt.shrt, true)
		}

		res.Body.Close()
	}
}

func TestShrtListHandler(t *testing.T) {
	shrts := []*goshrt.Shrt{
		{
			ID:     1,
			Domain: "domain",
			Slug:   "slug1",
			Dest:   "https://github.com/storvik/goshrt",
		},
		{
			ID:     2,
			Domain: "domain2",
			Slug:   "slug2",
			Dest:   "https://github.com/storvik/goshrt",
		},
		{
			ID:     3,
			Domain: "domain",
			Slug:   "slug3",
			Dest:   "https://github.com/storvik/goshrt",
		},
		{
			ID:     4,
			Domain: "domain2",
			Slug:   "slug4",
			Dest:   "https://github.com/storvik/goshrt",
		},
	}

	var tests = []struct {
		shrts              []*goshrt.Shrt
		expectedShrts      []*goshrt.Shrt
		url                string
		expectedStatusCode int
	}{
		{
			shrts:              shrts,
			expectedShrts:      shrts,
			url:                "/api/shrts",
			expectedStatusCode: http.StatusOK,
		},
		{
			shrts:              shrts,
			expectedShrts:      []*goshrt.Shrt{shrts[0], shrts[2]},
			url:                "/api/shrts/domain",
			expectedStatusCode: http.StatusOK,
		},
		{
			shrts:              shrts,
			url:                "/api/shrts/nonexistingdomain",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			shrts:              shrts,
			url:                "/api/shrts/error", // mockStore returns error if domain is error
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		// Create mock store and populate it with shrt.
		m := &mockStore{shrts: tt.shrts}

		s := newMockServer(t, m)
		req := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)

		w := executeRequest(req, s)

		res := w.Result()

		checkResponseCode(t, tt.expectedStatusCode, res.StatusCode)

		// Check if received shrt is correct, but only if http status is ok.
		if res.StatusCode == http.StatusOK {
			var resShrts []*goshrt.Shrt
			decoder := json.NewDecoder(res.Body)

			err := decoder.Decode(&resShrts)
			if err != nil {
				t.Errorf("error decoding response, %v", err)
			}

			for c, r := range resShrts {
				compareShrts(t, r, tt.expectedShrts[c], true)
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
		{
			shrt:               &goshrt.Shrt{},
			url:                "/api/shrt/3",
			expectedStatusCode: http.StatusInternalServerError,
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
			compareShrts(t, resShrt, tt.shrt, true)
		}

		res.Body.Close()
	}
}
