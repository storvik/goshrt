package postgres_test

import (
	"errors"
	"testing"
	"time"

	"github.com/storvik/goshrt"
)

func TestShrtStorerPostgres_CreateShrt(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s := &goshrt.Shrt{
			Domain: "localhost:3000",
			Slug:   "test",
			Dest:   "https://github.com/storvik/goshrt",
		}

		if err := db.CreateShrt(s); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ErrMultiple", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		// TODO: Add better tests that tests various timestamps
		s := &goshrt.Shrt{
			Domain: "gotest.com",
			Slug:   "TestMultiple",
			Dest:   "http://github.com/storvik/goshrt",
			Expiry: time.Now().Add(24 * time.Hour),
		}

		if err := db.CreateShrt(s); err != nil {
			t.Fatal(err)
		}

		err := db.CreateShrt(s)
		if err == nil {
			t.Fatal("expected error, but received none")
		}
		if !errors.Is(err, goshrt.ErrMultiple) {
			t.Error("expected multiple error, received another error")
		}
	})

	t.Run("ErrInvalid", func(t *testing.T) {
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
			err := db.CreateShrt(tt)
			if err == nil {
				t.Fatal("expected error, but received none")
			}
			if !errors.Is(err, goshrt.ErrInvalid) {
				t.Error("expected invalid input error, received another error")
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

func TestShrtStorerPostgres_ShrtByID(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s1 := &goshrt.Shrt{
			Domain: "gotest.com",
			Slug:   "TestShrtStorerPostgres_ShrtByID",
			Dest:   "http://github.com/storvik/goshrt",
		}

		err := db.CreateShrt(s1)
		if err != nil {
			t.Fatal(err)
		}

		s2, err := db.ShrtByID(s1.ID)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Should probably test timestamp too
		if s1.Domain != s2.Domain || s1.Slug != s2.Slug || s1.Dest != s2.Dest {
			t.Error("input not equal to output")
		}
	})
}

func TestShrtStorerPostgres_DeleteByID(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		s1 := &goshrt.Shrt{
			Domain: "gotest.com",
			Slug:   "TestShrtStorerPostgres_DeleteByID",
			Dest:   "http://github.com/storvik/goshrt",
		}

		err := db.CreateShrt(s1)
		if err != nil {
			t.Fatal(err)
		}

		s2, err := db.DeleteByID(s1.ID)
		if err != nil {
			t.Fatal(err)
		}

		if s1.Domain != s2.Domain || s1.Slug != s2.Slug || s1.Dest != s2.Dest {
			t.Error("input not equal to deleted shrt")
		}

		_, err = db.ShrtByID(s1.ID)
		if !errors.Is(err, goshrt.ErrNotFound) {
			t.Error("expected to not find deleted shrt, but found it")
		}
	})
}

func TestShrtStorerPostgres_Shrts(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		_, err := db.Shrts()
		if err != nil {
			t.Fatal(err)
		}
		// TODO: Find a better test case here
	})
}

func TestShrtStorerPostgres_ShrtsByDomain(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t)
		defer MustCloseDB(t, db)

		var tests = []*goshrt.Shrt{
			{
				Domain: "shrtsbydomain.com",
				Slug:   "1",
				Dest:   "http://github.com/storvik/goshrt",
			},
			{
				Domain: "shrtsbydomain.com",
				Slug:   "2",
				Dest:   "http://github.com/storvik/goshrt",
			},
			{
				Domain: "shrtsbydomain.com",
				Slug:   "3",
				Dest:   "http://github.com/storvik/goshrt",
			},
			{
				Domain: "shrtsbydomain.com",
				Slug:   "4",
				Dest:   "http://github.com/storvik/goshrt",
			},
		}

		for _, tt := range tests {
			err := db.CreateShrt(tt)
			if err != nil {
				t.Fatal(err)
			}
		}

		shrts, err := db.ShrtsByDomain("shrtsbydomain.com")
		if err != nil {
			t.Fatal(err)
		}
		if len(shrts) != len(tests) {
			t.Errorf("expected %d items with given domain, got %d", len(tests), len(shrts))
		}
	})
}
