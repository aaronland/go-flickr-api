package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/response"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

const API_ENDPOINT string = "https://api.flickr.com/services/rest"
const UPLOAD_ENDPOINT string = "https://up.flickr.com/services/upload/"
const REPLACE_ENDPOINT string = "https://up.flickr.com/services/replace/"

type Client interface {
	WithAccessToken(context.Context, auth.AccessToken) (Client, error)
	GetRequestToken(context.Context, string) (auth.RequestToken, error)
	GetAuthorizationURL(context.Context, auth.RequestToken, string) (string, error)
	GetAccessToken(context.Context, auth.RequestToken, auth.AuthorizationToken) (auth.AccessToken, error)
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
	Upload(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	Replace(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
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

func UploadAsyncWithClient(ctx context.Context, cl Client, fh io.Reader, args *url.Values) (int64, error) {

	args.Set("async", "1")

	rsp, err := cl.Upload(ctx, fh, args)

	if err != nil {
		return 0, err
	}

	return checkAsyncResponseWithClient(ctx, cl, rsp)
}

func ReplaceAsyncWithClient(ctx context.Context, cl Client, fh io.Reader, args *url.Values) (int64, error) {

	args.Set("async", "1")

	rsp, err := cl.Replace(ctx, fh, args)

	if err != nil {
		return 0, err
	}

	return checkAsyncResponseWithClient(ctx, cl, rsp)
}

func checkAsyncResponseWithClient(ctx context.Context, cl Client, rsp_fh io.ReadSeekCloser) (int64, error) {

	ticket, err := response.UnmarshalTicketResponse(rsp_fh)

	if err != nil {
		return 0, err
	}

	if ticket.Error != nil {
		return 0, ticket.Error
	}

	if ticket.TicketId == "" {
		return 0, fmt.Errorf("Missing ticket ID")
	}

	return CheckTicketWithClient(ctx, cl, ticket)
}

func CheckTicketWithClient(ctx context.Context, cl Client, ticket *response.Ticket) (int64, error) {

	// SET TIMEOUT HERE

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, nil
		case <-ticker.C:

			args := &url.Values{}
			args.Set("method", "flickr.photos.upload.checkTickets")
			args.Set("tickets", ticket.TicketId)
			// args.Set("format", "rest")

			// {"uploader":{"ticket":[{"id":"161192644-72157718847464617","complete":1,"photoid":"51090628667","imported":"1617400985"}]},"stat":"ok"}

			check_rsp, err := cl.ExecuteMethod(ctx, args)

			if err != nil {
				return 0, err
			}

			io.Copy(os.Stdout, check_rsp)
		}
	}

	return 0, nil
}
