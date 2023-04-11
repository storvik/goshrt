package goshrt_test

import (
	"testing"

	"github.com/storvik/goshrt"
)

func TestShrt_Validation(t *testing.T) {

	t.Run("Validate_Destination", func(t *testing.T) {

		var tests = []struct {
			s     goshrt.Shrt
			valid bool
		}{
			{goshrt.Shrt{Dest: "http://golang.org"}, true},
			{goshrt.Shrt{Dest: "http://golang.org:3000"}, true},
			{goshrt.Shrt{Dest: "https://golang.org"}, true},
			{goshrt.Shrt{Dest: "golang.org"}, false},
			{goshrt.Shrt{Dest: "https://golang.org/go"}, true},
			{goshrt.Shrt{Dest: "http://golang.org/go.html"}, true},
		}

		for _, tt := range tests {
			if v := tt.s.ValidDest(); v != tt.valid {
				t.Errorf("%s: expected valid: %t, got %t", tt.s.Dest, tt.valid, v)
			}
		}
	})

}
