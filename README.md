# go-flickr-api

Go package for working with the Flickr API

## Important

This is nearly, but not quite, complete. Some things may still change and not all the documentation is finished.

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

This package only defines [a handful of Go types or structs mapping to individual API responses](response). So far these are all specific to operations relating to uploading or replacing photos and to pagination.

In time there may be, along with helper methods for unmarshaling API responses in to typed responses but the baseline for all operations will remain: Query (`url.Values`) parameters sent over HTTP returning an `io.ReadSeekCloser` instance that is inspected and validated according to the needs and uses of the tools using the Flickr API.

## Tools

This package comes with a series of opinionated applications to implement functionality exposed by the Flickr API. These easiest way to build them is to run the handy `cli` target in the Makefile that comes bundled with this package.

_As of this writing the documentation and final output of all but the `api` tool is incomplete and subject to change still._

```
$> make cli
go build -mod vendor -o bin/api cmd/api/main.go
go build -mod vendor -o bin/upload cmd/upload/main.go
go build -mod vendor -o bin/replace cmd/replace/main.go
go build -mod vendor -o bin/auth-cli cmd/auth-cli/main.go
go build -mod vendor -o bin/auth-www cmd/auth-www/main.go
```

### api

Command-line tool for invoking the Flickr API. Results are emitted to STDOUT. Uploading and replacing images are not supported by this tool.

```
$> ./bin/api -h
Command-line tool for invoking the Flickr API. Results are emitted to STDOUT. Uploading and replacing images are not supported by this tool.

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

Command-line tool for initiating a Flickr API authorization flow.

```
> ./bin/auth-cli -h
Command-line tool for initiating a Flickr API authorization flow.

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

HTTP server for initiating a Flickr API autorization flow in a web browser.

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

Command-line tool for uploading an image to Flickr.

```
$> ./bin/upload -h
Usage of ./bin/upload:
Command-line tool for uploading an image to Flickr.

  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -param value
    	Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

### replace

Command-line tool for replacing an image in Flickr.

```
$> ./bin/replace -h
Command-line tool for replacing an image in Flickr.

Usage of ./bin/replace:
  -client-uri string
    	A valid aaronland/go-flickr-api client URI.
  -param value
    	Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.
  -use-runtimevar
    	Signal that all -uri flags are encoded as gocloud.dev/runtimevar string URIs.
```

### Design

The guts of all the tools bundled with this package are kept in the [application](application) directory rather than in application code itself. That's because the tools rely on the [GoCloud](https://gocloud.dev/) APIs for specific functionality. These are:

* Reading sensitive configuration data using the [runtimevar](https://gocloud.dev/howto/runtimevar/) abstraction layer.
* Reading images to upload or replace using the [blob](https://gocloud.dev/howto/blob/) abstraction layer.
* Persisting client authorization details (like request tokens) between stateless HTTP requests using the [docstore](https://gocloud.dev/howto/docstore/) abstraction layer.

Rather than bundling the code for all the services that the `GoCloud` APIs support the core of the application code only knows to expect _implementations_ of the relevant interfaces. It is expected that the application code itself (the code defined in `cmd/SOMEAPPLICATION`) will import the necessary packages for supporting a service-specific implementation.

For example, here's what the code for the [cmd/api](cmd/api/main.go) tool looks like:

```
package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/api"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	ctx := context.Background()

	app := &api.APIApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
```

This version of the application supports reading client configuration using the `runtimevar` package's [constant://](https://gocloud.dev/howto/runtimevar/#local) and [file://](https://gocloud.dev/howto/runtimevar/#local) schemes. For example:

```
$> bin/api -use-runtimevar -client-uri file:///path/to/client-uri.cfg -param method=flickr.test.echo
```

If you wanted a version of the `api` tool that supported reading client configuration stored in the Amazon Web Service's [Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html) secrets manager you would rewrite the application like this:

```
package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/api"
	_ "gocloud.dev/runtimevar/awsparamstore"
	"log"
)

func main() {

	ctx := context.Background()

	app := &api.APIApplication{}
	_, err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
```

And then invoke it like this:

```
$> bin/api -use-runtimevar -client-uri 'awsparamstore://{NAME}?region={REGION}&decoder=string' -param method=flickr.test.echo
```

The only thing that changes is the `runtimevar` package that your code imports. The rest of the application is encapsulated in the `APIApplication` instance.

## See also

* https://www.flickr.com/services/api/
* https://www.flickr.com/services/api/auth.oauth.html
* https://github.com/aaronland/go-http-server
* https://gocloud.dev/howto/docstore/
* https://gocloud.dev/howto/runtimevar/
* https://gocloud.dev/howto/blob/