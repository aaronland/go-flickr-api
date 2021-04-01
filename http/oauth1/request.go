package oauth1

import (
	"github.com/aaronland/go-flickr-api/client"
	"gocloud.dev/docstore"
	gohttp "net/http"
)

func NewRequestTokenHandler(cl client.Client, col *docstore.Collection, perms string) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()

		req_token, err := cl.GetRequestToken(ctx, req.Host) // FIX ME, callback URL

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		cache, err := NewRequestTokenCache(req_token)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		err = col.Put(ctx, cache)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		auth_url, err := cl.GetAuthorizationURL(ctx, req_token, perms)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		gohttp.Redirect(rsp, req, auth_url, gohttp.StatusFound)
	}

	return gohttp.HandlerFunc(fn), nil
}
