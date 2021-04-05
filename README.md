# go-flickr-api

Go package for working with the Flickr API

## Important

Work in progress. There may still be bugs. Complete documentation to follow.

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

This package only defines a handful of Go types or structs mapping to individual API responses. So far these are all specific to operations relating to uploading or replacing photos.

In time there may be, along with helper methods for unmarshaling API responses in to typed responses but the baseline for all operations will remain: Query (`url.Values`) parameters sent over HTTP returning an `io.ReadSeekCloser` instance that is inspected and validated according to the needs and uses of the tools using the Flickr API.

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
	WithAccessToken(context.Context, auth.AccessToken) (Client, error)	
	GetRequestToken(context.Context, string) (auth.RequestToken, error)
	GetAuthorizationURL(context.Context, auth.RequestToken, string) (string, error)
	GetAccessToken(context.Context, auth.RequestToken, auth.AuthorizationToken) (auth.AccessToken, error)
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
	Upload(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
	Replace(context.Context, io.Reader, *url.Values) (io.ReadSeekCloser, error)
}
```

### client.ExecuteMethodPaginatedCallback

```
type ExecuteMethodPaginatedCallback func(context.Context, io.ReadSeekCloser, error) error
```

## Tools

This package comes with a series of opinionated applications to implement functionality exposed by the Flickr API.

```
> make cli
go build -mod vendor -o bin/api cmd/api/main.go
go build -mod vendor -o bin/upload cmd/upload/main.go
go build -mod vendor -o bin/replace cmd/replace/main.go
go build -mod vendor -o bin/auth-cli cmd/auth-cli/main.go
go build -mod vendor -o bin/auth-www cmd/auth-www/main.go
```

### api

```
$> ./bin/api -h
Usage of ./bin/api:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -paginated
    	Automatically paginate (and iterate through) all API responses.
  -param value
    	Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

For example:

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

### auth-cli

```
> ./bin/auth-cli -h
Usage of ./bin/auth-cli:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -permissions string
    	A valid Flickr API permissions flag.
  -server-uri string
    	A valid aaronland/go-http-server URI.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

For example:

```
$> ./bin/auth-cli \
	-server-uri 'mkcert://localhost:8080' \
	-client-uri 'oauth1://?consumer_key={KEY}&consumer_secret={SECRET}'
	
2021/03/31 22:47:08 Checking whether mkcert is installed. If it is not you may be prompted for your password (in order to install certificate files
2021/03/31 22:47:09 Listening for requests on https://localhost:8080
2021/03/31 22:47:13 Authorize this application https://www.flickr.com/services/oauth/authorize?oauth_token={TOKEN}&perms=read

{"oauth_token":"{TOKEN}","oauth_token_secret":"{SECRET}"}
```

### auth-www

```
$> ./bin/auth-www -h
Usage of ./bin/auth-www:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -collection-uri string
    	A valid gocloud.dev/docstore URI. The docstore is used to store token requests during the time a user is approving an authentication request.
  -permissions string
    	A valid Flickr API permissions flag.
  -server-uri string
    	A valid aaronland/go-http-server URI.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

### upload

```
$> ./bin/upload -h
Usage of ./bin/upload:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -param value
    	Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

### replace

```
$> ./bin/replace -h
Usage of ./bin/replace:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -param value
    	Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

## See also

* https://www.flickr.com/services/api/
* https://www.flickr.com/services/api/auth.oauth.html
* https://github.com/aaronland/go-http-server
* https://gocloud.dev/howto/docstore/
* https://gocloud.dev/howto/runtimevar/
* https://gocloud.dev/howto/blob/