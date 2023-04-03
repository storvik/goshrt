package postgres

import (
	"database/sql"
	"fmt"

	"github.com/storvik/goshrt"
)

var _ goshrt.ShrtStorer = &client{}

type client struct {
	db *sql.DB

	name     string
	user     string
	password string
	address  string
}

func NewClient(n, u, p, a string) goshrt.ShrtStorer {
	return &client{
		name:     n,
		user:     u,
		password: p,
		address:  a,
	}
}

// Open connects to database using info stored in client.
func (c *client) Open() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", c.user, c.password, c.address, c.name)
	var err error
	c.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return nil
}

// Close closes database connection.
func (c *client) Close() error {
	return c.db.Close()
}

// Migrate aims to migrate database, gracefully.
func (c *client) Migrate() error {
	m := `
 CREATE TABLE IF NOT EXISTS shrts(
   domain text,     -- NOT NULL due to PK below
   slug text,       -- NOT NULL due to PK below
   dest text,
   expiry date,
   PRIMARY KEY (domain, slug)
 );`
	_, err := c.db.Exec(m)
	return err
}

// Shrt tries to find url in database
func (c *client) Shrt(d, s string) (*goshrt.Shrt, error) {
	shrt := &goshrt.Shrt{
		Domain: d,
		Slug:   s,
	}
	err := c.db.QueryRow("SELECT dest, expiry FROM ? WHERE domain=? AND slug", c.name, d, s).Scan(&shrt.Dest, &shrt.Expiry)
	if err == sql.ErrNoRows {
		// TODO: Add better error handling her, maybe custom domain type error
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return shrt, nil
}

// CreateShrt creates new shrt in database
func (c *client) CreateShrt(s *goshrt.Shrt) error {
	stmt, err := c.db.Prepare("INSERT INTO ? (domain, slug, dets, expiry) VALUES( ?, ?, ?, ? )")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(c.name, s.Domain, s.Slug, s.Dest, s.Expiry); err != nil {
		return err
	}
	return nil
}
