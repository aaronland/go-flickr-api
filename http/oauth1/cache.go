package oauth1

import (
	"github.com/aaronland/go-flickr-api/auth"
	"time"
)

// RequestTokenCache is a struct containing OAuth1 request token details and a timestamp
// indicating when the token details were created. This information is used to persist
// request token information, specifically the request token secret, before and after the
// approval phase of an OAuth1 authorization "www" flow.
type RequestTokenCache struct {
	// Flickr API OAuth1 request token.
	Token   string
	// Flickr API OAuth1 request token secret.	
	Secret  string
	// Unix timestamp representing the time that the RequestTokenCache was created.
	Created int64
}

// Create a new RequestTokenCache instance from a auth.RequestToken instance.
func NewRequestTokenCache(req_token auth.RequestToken) (*RequestTokenCache, error) {

	now := time.Now()

	cache := &RequestTokenCache{
		Token:   req_token.Token(),
		Secret:  req_token.Secret(),
		Created: now.Unix(),
	}

	return cache, nil
}
