package oauth1

import (
	"github.com/aaronland/go-flickr-api/client"
	"gocloud.dev/docstore"
	_ "log"
	gohttp "net/http"
	_ "net/url"
)

type RequestTokenHandlerOptions struct {
	Client       client.Client
	Collection   *docstore.Collection
	Permissions  string
	AuthCallback string
}

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
