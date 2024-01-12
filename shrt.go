package goshrt

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	slugAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Shrt struct {
	ID     int       `json:"id"`          // ID
	Domain string    `json:"domain"`      // Domain
	Slug   string    `json:"slug"`        // Everything avfter domain
	Dest   string    `json:"destination"` // Destination URL
	Expiry time.Time `json:"expire"`      // Timestamp for expire
}

// Printp pretty prints Shrt struct.
func (s *Shrt) Printp() {
	var timestring string
	if s.Expiry.IsZero() {
		timestring = "not set"
	} else {
		timestring = s.Expiry.Format("2006.02.01")
	}

	fmt.Printf("ID\t\t%d\n", s.ID)
	fmt.Printf("Domain\t\t%s\n", s.Domain)
	fmt.Printf("Slug\t\t%s\n", s.Slug)
	fmt.Printf("Destination\t%s\n", s.Dest)
	fmt.Printf("Expiry\t\t%s\n", timestring)
}

func (s *Shrt) ValidDest() bool {
	_, err := url.ParseRequestURI(s.Dest)
	if err != nil {
		return false
	}

	u, err := url.Parse(s.Dest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

type ShrtStorer interface {
	Migrate() error
	Open() error
	Close() error
	CreateShrt(s *Shrt) error
	Shrt(d, s string) (*Shrt, error)
	ShrtByID(id int) (*Shrt, error)
	DeleteByID(id int) (*Shrt, error)
	Shrts() ([]*Shrt, error)
	ShrtsByDomain(d string) ([]*Shrt, error)
	// TODO: Same as above with pagination
}

// GenerateSlug generates slug with length l.
func GenerateSlug(l uint64) string {
	var encodedBuilder strings.Builder

	length := len(slugAlphabet)

	encodedBuilder.Grow(10)

	for ; l > 0; l /= uint64(length) {
		encodedBuilder.WriteByte(slugAlphabet[(l % uint64(length))])
	}

	return encodedBuilder.String()
}

// ValidateSlug validates slug.
func ValidateSlug(s string) bool {
	r := regexp.MustCompile(`public(\/|$)|api(\/|$)|\?|#`)
	return !r.MatchString(s)
}
