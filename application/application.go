package application

import (
	"context"
	"flag"
)

// The Application interface is a common interface for the tools bundled with this package.
type Application interface {
	DefaultFlagSet() *flag.FlagSet
	Run(context.Context) (interface{}, error)
	RunWithFlagSet(context.Context, *flag.FlagSet) (interface{}, error)
}
