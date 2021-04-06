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
