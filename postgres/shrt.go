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
	err := c.db.QueryRow("SELECT dest, expiry FROM shrts WHERE domain=$1 AND slug=$2", d, s).Scan(&shrt.Dest, &shrt.Expiry)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return shrt, nil
}

// Shrt tries to find url in database by id
func (c *client) ShrtByID(id int) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{}
	err := c.db.QueryRow("SELECT domain, slug, dest, expiry FROM shrts WHERE id=$1", id).Scan(&shrt.Domain, &shrt.Slug, &shrt.Dest, &shrt.Expiry)
	if err == sql.ErrNoRows {
		return nil, goshrt.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return shrt, nil
}

// CreateShrt creates new shrt in database and returns its id. Has to look for
// same domain and slug which isn't expired.
func (c *client) CreateShrt(s *goshrt.Shrt) error {
	if s.Dest == "" || s.Domain == "" || s.Slug == "" {
		return goshrt.ErrInvalid
	}
	var e time.Time
	err := c.db.QueryRow("SELECT expiry FROM shrts WHERE domain=$1 AND slug=$2", s.Domain, s.Slug).Scan(&e)
	if err != sql.ErrNoRows {
		// TODO: This does not check for multiple expired rows
		if time.Now().Before(e) {
			return goshrt.ErrMultiple
		}
	}

	stmt, err := c.db.Prepare("INSERT INTO shrts(domain, slug, dest, expiry) VALUES( $1, $2, $3, $4 ) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(s.Domain, s.Slug, s.Dest, s.Expiry).Scan(&id)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}
