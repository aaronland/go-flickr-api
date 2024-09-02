// package application provides series of opinionated applications to implement functionality exposed by the Flickr API.
//
// The guts of all the tools bundled with this package are kept in this directory rather than in application code itself. That's because the tools rely on the GoCloud APIs for specific functionality. These are:
//
// * Reading sensitive configuration data using the runtimevar abstraction layer.
// * Reading images to upload or replace using the blob abstraction layer.
// * Persisting client authorization details (like request tokens) between stateless HTTP requests using the docstore abstraction layer.
//
// Rather than bundling the code for all the services that the `GoCloud` APIs support the core of the application code only knows to expect _implementations_ of the relevant interfaces. It is expected that the application code itself (the code defined in `cmd/SOMEAPPLICATION`) will import the necessary packages for supporting a service-specific implementation.
//
// For example, here's what the code for the cmd/api/main.go tool looks like:
//
// ```
// package main
//
// import (
//
//	"context"
//	"github.com/aaronland/go-flickr-api/application/api"
//	_ "gocloud.dev/runtimevar/constantvar"
//	_ "gocloud.dev/runtimevar/filevar"
//	"log"
//
// )
//
// func main() {
//
//		ctx := context.Background()
//
//		app := &api.APIApplication{}
//		_, err := app.Run(ctx)
//
//		if err != nil {
//			log.Fatalf("Failed to run upload application, %v", err)
//		}
//	}
//
// ```
//
// This version of the application supports reading client configuration using the `runtimevar` package's constant:// and file:// schemes. For example:
//
// ```
// $> bin/api -use-runtimevar -client-uri file:///path/to/client-uri.cfg -param method=flickr.test.echo
// ```
//
// If you wanted a version of the `api` tool that supported reading client configuration stored in the Amazon Web Service's Parameter Store secrets manager you would rewrite the application like this:
//
// ```
// package main
//
// import (
//
//	"context"
//	"github.com/aaronland/go-flickr-api/application/api"
//	_ "gocloud.dev/runtimevar/awsparamstore"
//	"log"
//
// )
//
// func main() {
//
//		ctx := context.Background()
//
//		app := &api.APIApplication{}
//		_, err := app.Run(ctx)
//
//		if err != nil {
//			log.Fatalf("Failed to run upload application, %v", err)
//		}
//	}
//
// ```
//
// And then invoke it like this:
//
// ```
// $> bin/api -use-runtimevar -client-uri 'awsparamstore://{NAME}?region={REGION}&decoder=string' -param method=flickr.test.echo
// ```
//
// The only thing that changes is the `runtimevar` package that your code imports. The rest of the application is encapsulated in the `APIApplication` instance.
package application

import (
	"context"
	"flag"
)

// The Application interface is a common interface for the tools bundled with this package.
type Application interface {
	// Return the default FlagSet necessary for the application to run.
	DefaultFlagSet() *flag.FlagSet
	// Invoke the application with its default FlagSet.
	Run(context.Context) (interface{}, error)
	// Invoke the application with a custom FlagSet.
	RunWithFlagSet(context.Context, *flag.FlagSet) (interface{}, error)
}
