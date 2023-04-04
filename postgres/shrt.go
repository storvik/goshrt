package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/storvik/goshrt"
)

var _ goshrt.ShrtStorer = &client{}

// Shrt tries to find url in database
func (c *client) Shrt(d, s string) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{
		Domain: d,
		Slug:   s,
	}
	err := c.db.QueryRow("SELECT dest, expiry FROM shrts WHERE domain=$1 AND slug=$2", d, s).Scan(&shrt.Dest, &shrt.Expiry)
	if err == sql.ErrNoRows {
		// TODO: Add better error handling here, maybe custom domain type error
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return shrt, nil
}

// CreateShrt creates new shrt in database
func (c *client) CreateShrt(s *goshrt.Shrt) error {
	if s.Dest == "" || s.Domain == "" || s.Slug == "" {
		return goshrt.ErrInvalid
	}
	stmt, err := c.db.Prepare("INSERT INTO shrts(domain, slug, dest, expiry) VALUES( $1, $2, $3, $4 )")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(s.Domain, s.Slug, s.Dest, s.Expiry); err != nil {
		return err
	}
	return nil
}
