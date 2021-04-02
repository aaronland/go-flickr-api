package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api"
	"github.com/aaronland/go-flickr-api/auth"
	"io"
	"net/url"
	"strconv"
)

const API_ENDPOINT string = "https://api.flickr.com/services/rest"
const UPLOAD_ENDPOINT string = "https://up.flickr.com/services/upload/"

type Client interface {
	WithAccessToken(context.Context, auth.AccessToken) (Client, error)
	GetRequestToken(context.Context, string) (auth.RequestToken, error)
	GetAuthorizationURL(context.Context, auth.RequestToken, string) (string, error)
	GetAccessToken(context.Context, auth.RequestToken, auth.AuthorizationToken) (auth.AccessToken, error)
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
	Upload(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	UploadAsync(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	Replace(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	ReplaceAsync(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
}

type ExecuteMethodPaginatedCallback func(context.Context, io.ReadSeekCloser, error) error

func ExecuteMethodPaginated(ctx context.Context, cl Client, args *url.Values, cb ExecuteMethodPaginatedCallback) error {

	page := 1
	pages := -1

	if args.Get("page") == "" {
		args.Set("page", strconv.Itoa(page))
	} else {

		p, err := strconv.Atoi(args.Get("page"))

		if err != nil {
			return fmt.Errorf("Invalid page number '%s', %v", args.Get("page"), err)
		}

		page = p
	}

	for {

		fh, err := cl.ExecuteMethod(ctx, args)

		err = cb(ctx, fh, err)

		if err != nil {
			return err
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to rewind response, %v", err)
		}

		if pages == -1 {

			pagination, err := api.DerivePagination(ctx, fh)

			if err != nil {
				return err
			}

			pages = pagination.Pages
		}

		page += 1

		if page <= pages {
			args.Set("page", strconv.Itoa(page))
		} else {
			break
		}
	}

	return nil
}
