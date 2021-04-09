package oauth1

import (
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-http-sanitize"
	gohttp "net/http"
)

// Return a new HTTP handler to receive a process OAuth1 authorization callback requests. This handler will
// relay the OAuth1 authorization token or any errors received by the callback to the appropriate channel.
// This handler is used to create a background HTTP server process that can block execution of a command-line
// OAuth1 authorization "www" flow until either a token or an error is dispatched to its corresponding channel
// in the application code.
func NewAuthorizationTokenHandlerWithChannels(token_ch chan auth.AuthorizationToken, err_ch chan error) (gohttp.Handler, error) {

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

		auth_token := &auth.OAuth1AuthorizationToken{
			OAuthToken:    token,
			OAuthVerifier: verifier,
		}

		token_ch <- auth_token

		rsp.Write([]byte(`Authorization request successful. You can close this browser window and return to the authorization application.`))
		return
	}

	return gohttp.HandlerFunc(fn), nil
}
