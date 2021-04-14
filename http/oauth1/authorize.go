package oauth1

import (
	"encoding/json"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-http-sanitize"
	"gocloud.dev/docstore"
	_ "io"
	"log"
	gohttp "net/http"
	"net/url"
	"os"
)

// AuthorizationTokenHandlerOptions is a struct containing application-specific details
// necessary for all OAuth1 authorization callback requests.
type AuthorizationTokenHandlerOptions struct {
	// A client.Client instance used to call the Flickr API
	Client client.Client
	// A gocloud.dev/docstore.Collection instance used to retrieve request token details necessary for creating permanent access tokens.
	Collection *docstore.Collection
}

// Return a new HTTP handler to receive a process OAuth1 authorization callback requests. This handler will
// retrieve the request token associated with the authorization request and exchange these elements for a permanent
// OAuth1 access token.
func NewAuthorizationTokenHandler(opts *AuthorizationTokenHandlerOptions) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()

		token, err := sanitize.GetString(req, "oauth_token")

		if err != nil {
			gohttp.Error(rsp, "Missing ?oauth_token parameter", gohttp.StatusBadRequest)
			return
		}

		verifier, err := sanitize.GetString(req, "oauth_verifier")

		if err != nil {
			gohttp.Error(rsp, "Missing ?oauth_verifier parameter", gohttp.StatusBadRequest)
			return
		}

		cache := &RequestTokenCache{
			Token: token,
		}

		err = opts.Collection.Get(ctx, cache)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		defer func() {

			err := opts.Collection.Delete(ctx, cache)

			if err != nil {
				log.Printf("Failed to delete cache item for %s, err\n", cache.Token, err)
			}
		}()

		req_token := &auth.OAuth1RequestToken{
			OAuthToken:       cache.Token,
			OAuthTokenSecret: cache.Secret,
		}

		auth_token := &auth.OAuth1AuthorizationToken{
			OAuthToken:    token,
			OAuthVerifier: verifier,
		}

		access_token, err := opts.Client.GetAccessToken(ctx, req_token, auth_token)

		if err != nil {
			gohttp.Error(rsp, "Missing ?oauth_verifier parameter", gohttp.StatusBadRequest)
			return
		}

		cl, err := opts.Client.WithAccessToken(ctx, access_token)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		args := &url.Values{}
		args.Set("method", "flickr.test.login")

		_, err = cl.ExecuteMethod(ctx, args)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		// STORE auth_token... WHERE? AND THEN WHAT?

		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(access_token)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Write([]byte(`Authorization request successful.`))
		return
	}

	return gohttp.HandlerFunc(fn), nil
}
