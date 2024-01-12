package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // pq, postgres sql driver
	"github.com/storvik/goshrt"
)

var _ goshrt.ShrtStorer = &Client{}

type Client struct {
	db *sql.DB

	name     string
	user     string
	password string
	address  string
	schema   string
}

func NewClient(n, u, p, a, s string) *Client {
	return &Client{
		name:     n,
		user:     u,
		password: p,
		address:  a,
		schema:   s,
	}
}

// Open connects to database using info stored in client.
func (c *Client) Open() error {
	var err error
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable&search_path=%s", c.user, c.password, c.address, c.name, c.schema)

	c.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes database connection.
func (c *Client) Close() error {
	return c.db.Close()
}

// Migrate aims to migrate database, gracefully.
func (c *Client) Migrate() error {
	m := `
 CREATE TABLE IF NOT EXISTS goshrt.shrts(
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
