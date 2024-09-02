// Example
// ```
// package main
//
// import (
//
//	"context"
//	"github.com/aaronland/go-flickr-api/client"
//	"io"
//	"net/url"
//	"os"
//
// )
//
// func main() {
//
//		ctx := context.Background()
//
//		client_uri := "oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth_token={TOKEN}&oauth_token_secret={SECRET}"
//		cl, _ := client.NewClient(ctx, client_uri)
//
//		args := &url.Values{}
//		args.Set("method", "flickr.test.login")
//
//		fh, _ := cl.ExecuteMethod(ctx, args)
//
//		defer fh.Close()
//
//		io.Copy(os.Stdout, fh)
//	}
//
// ```
//
// _Error handling removed for the sake of brevity._
//
// # Design
//
// The core of this package's approach to the Flickr API is the `ExecuteMethod` method (which is defined in the `client.Client` interface) whose signature looks like this:
//
// ```
//
//	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
//
// ```
//
// This package only defines [a handful of Go types or structs mapping to individual API responses](response). So far these are all specific to operations relating to uploading or replacing photos and to pagination.
//
// In time there may be, along with helper methods for unmarshaling API responses in to typed responses but the baseline for all operations will remain: Query (`url.Values`) parameters sent over HTTP returning an `io.ReadSeekCloser` instance that is inspected and validated according to the needs and uses of the tools using the Flickr API.
package api
