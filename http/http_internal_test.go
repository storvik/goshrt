package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/storvik/goshrt"
)

// executeRequest is a helper that executes request.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a helper that checks response code of request.
func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()

	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// compareShrts is a test helper that compares shrts.
func compareShrts(t *testing.T, shrt, expected *goshrt.Shrt, includeID bool) {
	t.Helper()

	eq := (shrt.Dest == expected.Dest)
	eq = eq && (shrt.Domain == expected.Domain)
	eq = eq && (shrt.Slug == expected.Slug)

	if includeID {
		eq = eq && (shrt.ID == expected.ID)
	}

	if !eq {
		t.Errorf("shrts not equal, expected %v, got %v", expected, shrt)
	}
}

// mockStore used in tests. All functions that returns
// one shrt will return shrt.
type mockStore struct {
	shrt  *goshrt.Shrt
	shrts []*goshrt.Shrt
}

func (m *mockStore) Migrate() error {
	return nil
}

func (m *mockStore) Open() error {
	return nil
}

func (m *mockStore) Close() error {
	return nil
}

func (m *mockStore) CreateShrt(s *goshrt.Shrt) error {
	if m.shrt.Domain == "error" {
		return goshrt.ErrInvalid
	}

	m.shrt = s

	return nil
}

// Shrt gets shrt from domain and slug.
func (m *mockStore) Shrt(d, s string) (*goshrt.Shrt, error) {
	if d == m.shrt.Domain && s == m.shrt.Slug {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

// ShrtByID gets shrt by ID.
func (m *mockStore) ShrtByID(id int) (*goshrt.Shrt, error) {
	// mock shrtstore not populated, used to test internal server error
	if m.shrt == nil || m.shrt.ID == 0 {
		return nil, errors.New("mock shrtstore not populated")
	}

	if id == m.shrt.ID {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

// DeleteByID deletes shrt by ID and returns deleted shrt.
func (m *mockStore) DeleteByID(id int) (*goshrt.Shrt, error) {
	// mock shrtstore not populated, used to test internal server error
	if m.shrt == nil || m.shrt.ID == 0 {
		return nil, errors.New("mock shrtstore not populated")
	}

	if id == m.shrt.ID {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

func (m *mockStore) Shrts() ([]*goshrt.Shrt, error) {
	return m.shrts, nil
}

func (m *mockStore) ShrtsByDomain(d string) ([]*goshrt.Shrt, error) {
	if d == "error" {
		return nil, errors.New("testerror")
	}

	shrts := []*goshrt.Shrt{}
	for _, s := range m.shrts {
		if s.Domain == d {
			shrts = append(shrts, s)
		}
	}

	if len(shrts) < 1 {
		return nil, goshrt.ErrNotFound
	}

	return shrts, nil
}

type mockAuthorizer struct {
}

func (m *mockAuthorizer) Create(_ string) (string, error) {
	return "", nil
}

func (m *mockAuthorizer) Validate(_ string) (bool, error) {
	return true, nil
}
