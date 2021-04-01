package oauth1

import (
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-http-sanitize"
	gohttp "net/http"
)

func NewAuthorizationTokenHandlerWithChannels(token_ch chan *auth.AuthorizationToken, err_ch chan error) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

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

		auth_token := &auth.AuthorizationToken{
			Token:    token,
			Verifier: verifier,
		}

		token_ch <- auth_token

		rsp.Write([]byte(`Authorization request successful.`))
		return
	}

	return gohttp.HandlerFunc(fn), nil
}
