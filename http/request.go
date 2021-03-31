package http

import (
	"github.com/aaronland/go-flickr-api/auth"
	"io"
	gohttp "net/http"
)

func NewOAuth1AuthorizeTokenHandler(token_ch chan *auth.AuthorizationToken, err_ch chan error) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		defer req.Body.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			err_ch <- err
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		auth_token, err := auth.UnmarshalAuthorizationToken(string(body))

		if err != nil {
			err_ch <- err
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		token_ch <- auth_token
		return
	}

	return gohttp.HandlerFunc(fn), nil
}
