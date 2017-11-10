package websocketemail

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Generate an email address of the form "^[a-f0-9]{32}\\@websocket\\.email$"
// using crypto/rand or return an error if it was not possible.
func GenerateEmailAddress() (string, error) {
	buf := make([]byte, 16, 16)
	// We can ignore the number of bytes read
	// for crypto/rand.Read(), it is not a reader.
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@websocket.email", hex.EncodeToString(buf)), nil
}

// Generate an email address of the form "^[a-f0-9]{32}\\@websocket\\.email$"
// using crypto/rand or panic if it was not possible.
func MustGenerateEmailAddress() string {
	email, err := GenerateEmailAddress()
	if err != nil {
		panic(err)
	}
	return email
}
