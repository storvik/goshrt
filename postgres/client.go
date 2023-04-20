package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
   id serial primary key,
   domain text not null,
   slug text not null,
   dest text not null,
   expiry date,
   deleted bool default false
);`
	_, err := c.db.Exec(m)
	return err
}
