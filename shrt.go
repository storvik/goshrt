package goshrt

import (
	"strings"
	"time"
)

const (
	slugAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Shrt struct {
	Domain string    `json:"domain"`      // Domain
	Slug   string    `json:"slug"`        // Everything avfter domain
	Dest   string    `json:"destination"` // Destination URL
	Expiry time.Time `json:"expire"`      // Timestamp for expire
}

type ShrtStorer interface {
	CreateShrt(s *Shrt) error
	Shrt(d, s string) (*Shrt, error)
	// TODO: Delete shrt and get list of shrts
}

// GenerateSlug generates slug with length l
func GenerateSlug(l uint64) string {
	length := len(slugAlphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; l > 0; l = l / uint64(length) {
		encodedBuilder.WriteByte(slugAlphabet[(l % uint64(length))])
	}
	return encodedBuilder.String()
}
