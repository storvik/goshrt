package postgres_test

import (
	"testing"

	"github.com/storvik/goshrt"
	"github.com/storvik/goshrt/postgres"
)

// Ensure the test database can open & close.
func TestDB(t *testing.T) {
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

// MustOpenDB creates a new database object, connects to it
// and runs migrations. Fatal error if any of the above fails.
// Needs a freshly setup database with:
//
//	Database name:      goshrt
//	Database user:      goshrt
//	Database password:  trhsog
//	Database address:   localhost:6000
func MustOpenDB(tb testing.TB) goshrt.ShrtStorer {
	tb.Helper() // This function is a test-helper, not a test

	db := postgres.NewClient("goshrt", "goshrt", "trhsog", "localhost:6000")
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		tb.Fatal(err)
	}

	return db
}

// MustCloseDB closes the DB, fatal on error.
func MustCloseDB(tb testing.TB, db goshrt.ShrtStorer) {
	tb.Helper()
	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}
