package oauth1

import (
	"github.com/aaronland/go-flickr-api/client"
	gohttp "net/http"
)

func NewRequestTokenHandler(cl client.Client, perms string) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()

		req_token, err := cl.GetRequestToken(ctx, req.Host) // FIX ME, callback URL

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		// STORE req_token.Secret in cache, key off req_token.Token

		auth_url, err := cl.GetAuthorizationURL(ctx, req_token, perms)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		gohttp.Redirect(rsp, req, auth_url, gohttp.StatusFound)
	}

	return gohttp.HandlerFunc(fn), nil
}
