package postgres_test

import (
	"testing"

	"github.com/storvik/goshrt"
)

func TestShrtStorerPostgres_CreateShrt(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s := &goshrt.Shrt{
			Domain: "gotest.com",
			Slug:   "123456",
			Dest:   "http://github.com/storvik/goshrt",
		}

		if err := db.CreateShrt(s); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ErrNameRequired", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		var tests = []*goshrt.Shrt{
			{
				Domain: "",
				Slug:   "1234",
				Dest:   "1234",
			},
			{
				Domain: "1234",
				Slug:   "",
				Dest:   "1234",
			},
			{
				Domain: "1234",
				Slug:   "1234",
				Dest:   "",
			},
		}

		for _, tt := range tests {
			if err := db.CreateShrt(tt); err == nil {
				t.Fatal("expected error, but received none")
			}
		}

	})
}

func TestShrtStorerPostgres_Shrt(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s1 := &goshrt.Shrt{
			Domain: "gotest.com",
			Slug:   "TestShrtStorerPostgres_Shrt",
			Dest:   "http://github.com/storvik/goshrt",
		}

		if err := db.CreateShrt(s1); err != nil {
			t.Fatal(err)
		}

		s2, err := db.Shrt(s1.Domain, s1.Slug)
		if err != nil {
			t.Fatal(err)
		}

		if s1.Domain != s2.Domain || s1.Slug != s2.Slug || s1.Dest != s2.Dest {
			t.Error("input not equal to output")
		}

	})

}
