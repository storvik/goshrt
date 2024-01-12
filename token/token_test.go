package token_test

import (
	"fmt"
	"testing"

	"github.com/storvik/goshrt/token"
)

func TestToken_Create(t *testing.T) {
	a := token.NewAuth("qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$")

	// Test that token can be created
	t.Run("OK", func(t *testing.T) {
		_, err := a.Create("testid")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestToken_Validate(t *testing.T) {
	a1 := token.NewAuth("qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$")

	valid, err := a1.Create("testid")
	if err != nil {
		t.Fatal(err)
	}

	// Slightly different from a1
	a2 := token.NewAuth("qTGVn$a&hLJ9385C^z7L!MW5CnwZq&$")

	invalid, err := a2.Create("testid")
	if err != nil {
		t.Fatal(err)
	}

	// Test if token can be validated
	t.Run("OK", func(t *testing.T) {
		var tests = []struct {
			token string
			valid bool
		}{
			{valid, true},
			{invalid, false},
		}

		for _, tt := range tests {
			testname := fmt.Sprintf("Token, valid: %t", tt.valid)
			t.Run(testname, func(t *testing.T) {
				p, _ := a1.Validate(tt.token)
				if p != tt.valid {
					t.Errorf("got %t, want %t", p, tt.valid)
				}
			})
		}
	})
}
