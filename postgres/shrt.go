package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/storvik/goshrt"
)

var _ goshrt.ShrtStorer = &client{}

// Shrt tries to find url in database by domain and slug
func (c *client) Shrt(d, s string) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{
		Domain: d,
		Slug:   s,
	}
	t := sql.NullTime{}
	err := c.db.QueryRow("SELECT id, dest, expiry FROM shrts WHERE domain=$1 AND slug=$2", d, s).Scan(&shrt.ID, &shrt.Dest, &t)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t.Valid {
		shrt.Expiry = t.Time
	}
	return shrt, nil
}

// Shrt tries to find url in database by id
func (c *client) ShrtByID(id int) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{
		ID: id,
	}
	t := sql.NullTime{}
	err := c.db.QueryRow("SELECT domain, slug, dest, expiry FROM shrts WHERE id=$1", id).Scan(&shrt.Domain, &shrt.Slug, &shrt.Dest, &t)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t.Valid {
		shrt.Expiry = t.Time
	}
	return shrt, nil
}

// CreateShrt creates new shrt in database and returns its id. Has to look for
// same domain and slug which isn't expired.
func (c *client) CreateShrt(s *goshrt.Shrt) error {
	if s.Dest == "" || s.Domain == "" || s.Slug == "" {
		return goshrt.ErrInvalid
	}
	rows, err := c.db.Query("SELECT expiry FROM shrts WHERE domain=$1 AND slug=$2", s.Domain, s.Slug)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var e time.Time
		err = rows.Scan(&e)
		if err != nil {
			return err
		}
		if time.Now().Before(e) || e.IsZero() {
			return goshrt.ErrMultiple
		}
	}
	expiry := sql.NullTime{
		Time:  s.Expiry,
		Valid: !s.Expiry.IsZero(),
	}
	stmt, err := c.db.Prepare("INSERT INTO shrts(domain, slug, dest, expiry) VALUES( $1, $2, $3, $4 ) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(s.Domain, s.Slug, s.Dest, expiry).Scan(&id)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

// DeleteByID deletes shrt by id
func (c *client) DeleteByID(id int) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{
		ID: id,
	}
	t := sql.NullTime{}
	err := c.db.QueryRow("SELECT domain, slug, dest, expiry FROM shrts WHERE id=$1", id).Scan(&shrt.Domain, &shrt.Slug, &shrt.Dest, &t)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t.Valid {
		shrt.Expiry = t.Time
	}
	_, err = c.db.Exec("DELETE FROM shrts WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return shrt, nil
}

// Shrts retrenves all shorts in database. Note that this is not efficient
// at all, should seriously add pagination to this.
func (c *client) Shrts() ([]*goshrt.Shrt, error) {
	var shrts []*goshrt.Shrt
	rows, err := c.db.Query("SELECT id, domain, slug, dest, expiry FROM shrts")
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &goshrt.Shrt{}
		t := sql.NullTime{}
		err = rows.Scan(&s.ID, &s.Domain, &s.Slug, &s.Dest, &t)
		if err != nil {
			return nil, err
		}
		if t.Valid {
			s.Expiry = t.Time
		}
		shrts = append(shrts, s)
	}

	return shrts, nil
}

// ShrtsByDomain retrieves all shrts by domain. Note that this is not efficient
// at all, should seriously add pagination to this.
func (c *client) ShrtsByDomain(d string) ([]*goshrt.Shrt, error) {
	var shrts []*goshrt.Shrt
	rows, err := c.db.Query("SELECT id, domain, slug, dest, expiry FROM shrts WHERE domain=$1", d)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		s := &goshrt.Shrt{}
		t := sql.NullTime{}
		err = rows.Scan(&s.ID, &s.Domain, &s.Slug, &s.Dest, &t)
		if err != nil {
			return nil, err
		}
		if t.Valid {
			s.Expiry = t.Time
		}
		shrts = append(shrts, s)
	}

	return shrts, nil
}
