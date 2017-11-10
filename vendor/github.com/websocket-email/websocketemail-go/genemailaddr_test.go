package websocketemail

import (
	"regexp"
	"testing"
)

func TestEmailAddressGeneration(t *testing.T) {
	for i := 0; i < 100; i++ {
		email := MustGenerateEmailAddress()
		ok, err := regexp.MatchString("^[a-f0-9]{32}\\@websocket\\.email$", email)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("invalid email: %s", email)
		}
	}
}
