package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/aaronland/go-flickr-api/response"
	"io"
	"net/url"
	"strconv"
	"time"
)

const API_ENDPOINT string = "https://api.flickr.com/services/rest"
const UPLOAD_ENDPOINT string = "https://up.flickr.com/services/upload/"
const REPLACE_ENDPOINT string = "https://up.flickr.com/services/replace/"

// Client is the interface that defines common methods for all Flickr API Client.
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

func ExecuteMethodPaginatedWithClient(ctx context.Context, cl Client, args *url.Values, cb ExecuteMethodPaginatedCallback) error {

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

			pagination, err := response.DerivePagination(ctx, fh)

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

	ticket, err := response.UnmarshalUploadTicketResponse(rsp_fh)

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

func CheckTicketWithClient(ctx context.Context, cl Client, ticket *response.UploadTicket) (int64, error) {

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

			check_rsp, err := cl.ExecuteMethod(ctx, args)

			if err != nil {
				return 0, err
			}

			check_ticket, err := response.UnmarshalCheckTicketResponse(check_rsp)

			if err != nil {
				return 0, err
			}

			for _, t := range check_ticket.Uploader.Tickets {

				if t.TicketId != ticket.TicketId {
					continue
				}

				if t.Complete != 1 {
					continue
				}

				// Because the Flickr API returns strings
				// for photo IDs

				str_id := t.PhotoId

				id, err := strconv.ParseInt(str_id, 10, 64)

				if err != nil {
					return 0, err
				}

				return id, nil
			}
		}
	}

	return 0, nil
}
