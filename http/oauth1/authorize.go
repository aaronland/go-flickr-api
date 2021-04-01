package oauth1

import (
	"encoding/json"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/aaronland/go-http-sanitize"
	"io"
	gohttp "net/http"
	"net/url"
)

func NewAuthorizationTokenHandler(cl client.Client) (gohttp.Handler, error) {

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

		auth_token := &auth.OAuth1AuthorizationToken{
			OAuthToken:    token,
			OAuthVerifier: verifier,
		}

		// RETRIEVE req_token FROM CACHE, key off auth_token.Token

		var req_token auth.RequestToken // FIX ME...
		access_token, err := cl.GetAccessToken(ctx, req_token, auth_token)

		if err != nil {
			gohttp.Error(rsp, "Missing ?oauth_verifier parameter", gohttp.StatusBadRequest)
			return
		}

		cl, err = cl.WithAccessToken(ctx, access_token)

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

		enc := json.NewEncoder(io.Discard)
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
