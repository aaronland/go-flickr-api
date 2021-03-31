package client

import (
	"context"
	"github.com/aaronland/go-flickr-api/auth"
	"io"
	"net/url"
)

type Client interface {
	GetRequestToken(context.Context, string) (*auth.RequestToken, error)
	AuthorizationURL(context.Context, *auth.RequestToken, string) (*url.URL, error)
	GetAccessToken(context.Context, *auth.AuthorizationToken) (*auth.AccessToken, error)
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
}
