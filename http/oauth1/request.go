package oauth1

import (
	"github.com/aaronland/go-flickr-api/client"
	"gocloud.dev/docstore"
	_ "log"
	gohttp "net/http"
)

// RequestTokenHandlerOptions is a struct containing application-specific details
// necessary for all OAuth1 authorization flow requests.
type RequestTokenHandlerOptions struct {
	// A client.Client instance used to call the Flickr API	
	Client       client.Client
	// A gocloud.dev/docstore.Collection instance used to store request token details necessary for creating permanent access tokens.	
	Collection   *docstore.Collection
	// The Flickr API permissions that your application is requesting.
	Permissions  string
	// The fully qualified callback URL to be invoked by Flickr if an autorization request is approved.
	AuthCallback string
}

// Return a new HTTP handler to create a new OAuth1 authorization request token and then redirect to the
// Flickr API OAuth1 authorization approval endpoint.
func NewRequestTokenHandler(opts *RequestTokenHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()

		req_token, err := opts.Client.GetRequestToken(ctx, opts.AuthCallback)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		cache, err := NewRequestTokenCache(req_token)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		err = opts.Collection.Put(ctx, cache)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		auth_url, err := opts.Client.GetAuthorizationURL(ctx, req_token, opts.Permissions)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		gohttp.Redirect(rsp, req, auth_url, gohttp.StatusFound)
	}

	return gohttp.HandlerFunc(fn), nil
}
