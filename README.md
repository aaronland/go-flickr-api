# go-flickr-api

Go package for working with the Flickr API

## Important

Work in progress. Uploads not supported yet. There may still be bugs. Complete documentation to follow.

## Example

```
package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/client"
	"io"
	"net/url"
	"os"
)

func main() {

	ctx := context.Background()

	client_uri := "oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth_token={TOKEN}&oauth_token_secret={SECRET}"
	cl, _ := client.NewClient(ctx, client_uri)

	args := &url.Values{}
	args.Set("method", "flickr.test.login")

	fh, _ := cl.ExecuteMethod(ctx, args)

	defer fh.Close()
	
	io.Copy(os.Stdout, fh)
}
```

_Error handling removed for the sake of brevity._

## Design

The core of this package's approach to the Flickr API is the `ExecuteMethod` method (which is defined in the `client.Client` interface) whose signature looks like this:

```
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
```

This package does not define any Go types or structs mapping to individual API responses yet. In time there may be, along with helper methods for unmarshaling API responses in to typed responses but the baseline for all operations will remain: Query (`url.Values`) parameters sent over HTTP returning an `io.ReadSeekCloser` instance that is inspected and validated according to the needs and uses of the tools using the Flickr API.

## Interfaces

### auth.RequestToken

```
type RequestToken interface {
	Token() string
	Secret() string
}
```

### auth.AuthorizationToken

```
type AuthorizationToken interface {
	Token() string
	Verifier() string
}
```

### auth.AccessToken

```
type AccessToken interface {
	Token() string
	Secret() string
}
```

### client.Client

```
type Client interface {
	GetRequestToken(context.Context, string) (auth.RequestToken, error)
	GetAuthorizationURL(context.Context, auth.RequestToken, string) (string, error)
	GetAccessToken(context.Context, auth.RequestToken, auth.AuthorizationToken) (auth.AccessToken, error)
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
	ExecuteMethodPaginated(context.Context, *url.Values, ExecuteMethodPaginatedCallback) error	
	WithAccessToken(context.Context, auth.AccessToken) (Client, error)
}
```

### client.ExecuteMethodPaginatedCallback

```
type ExecuteMethodPaginatedCallback func(context.Context, io.ReadSeekCloser, error) error
```

_Important: This interface may still change._

## Tools

```
$> make cli
go build -mod vendor -o bin/api cmd/api/main.go
go build -mod vendor -o bin/authorize cmd/authorize/main.go
```

### api

```
$> ./bin/api \
	-client-uri 'oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth_token={TOKEN}&oauth_token_secret={SECRET}' \
	-param method=flickr.test.login \

| jq

{
  "user": {
    "id": "123456789@X03",
    "username": {
      "_content": "example"
    },
    "path_alias": null
  },
  "stat": "ok"
}
```

## authorize

```
$> ./bin/authorize \
	-server-uri 'mkcert://localhost:8080' \
	-client-uri 'oauth1://?consumer_key={KEY}&consumer_secret={SECRET}'
	
2021/03/31 22:47:08 Checking whether mkcert is installed. If it is not you may be prompted for your password (in order to install certificate files
2021/03/31 22:47:09 Listening for requests on https://localhost:8080
2021/03/31 22:47:13 Authorize this application https://www.flickr.com/services/oauth/authorize?oauth_token={TOKEN}&perms=read

{"oauth_token":"{TOKEN}","oauth_token_secret":"{SECRET}"}
```

## See also

* https://www.flickr.com/services/api/
* https://www.flickr.com/services/api/auth.oauth.html