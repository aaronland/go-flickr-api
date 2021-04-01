package auth

import (
	"math/rand"
	"time"
)

type RequestToken interface {
	Token() string
	Secret() string
}

type AuthorizationToken interface {
	Token() string
	Verifier() string
}

type AccessToken interface {
	Token() string
	Secret() string
}

// Generate a random string of 8 chars, needed for OAuth signature
func GenerateNonce() string {

	rand.Seed(time.Now().UTC().UnixNano())

	// For convenience, use a set of chars we don't need to url-escape
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
