package oauth1

import (
	"github.com/aaronland/go-flickr-api/auth"
	"time"
)

type RequestTokenCache struct {
	Token   string
	Secret  string
	Created int64
}

func NewRequestTokenCache(req_token auth.RequestToken) (*RequestTokenCache, error) {

	now := time.Now()

	cache := &RequestTokenCache{
		Token:   req_token.Token(),
		Secret:  req_token.Secret(),
		Created: now.Unix(),
	}

	return cache, nil
}
