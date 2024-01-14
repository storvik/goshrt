package http

import (
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

// mockStore used in tests. All functions that returns
// one shrt will return shrt.
type mockStore struct {
	shrt *goshrt.Shrt
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

func (m *mockStore) CreateShrt(_ *goshrt.Shrt) error {
	return nil
}

// Shrt gets shrt from domain and slug.
func (m *mockStore) Shrt(d, s string) (*goshrt.Shrt, error) {
	if d == m.shrt.Domain && s == m.shrt.Domain {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

// ShrtByID gets shrt by ID.
func (m *mockStore) ShrtByID(id int) (*goshrt.Shrt, error) {
	if id == m.shrt.ID {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

// DeleteByID deletes shrt by ID and returns deleted shrt.
func (m *mockStore) DeleteByID(id int) (*goshrt.Shrt, error) {
	if id == m.shrt.ID {
		return m.shrt, nil
	}

	return nil, goshrt.ErrNotFound
}

func (m *mockStore) Shrts() ([]*goshrt.Shrt, error) {
	return nil, nil
}

func (m *mockStore) ShrtsByDomain(_ string) ([]*goshrt.Shrt, error) {
	return nil, nil
}

type mockAuthorizer struct {
}

func (m *mockAuthorizer) Create(_ string) (string, error) {
	return "", nil
}

func (m *mockAuthorizer) Validate(_ string) (bool, error) {
	return true, nil
}
