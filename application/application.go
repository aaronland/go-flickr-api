package application

import (
	"context"
	"flag"
)

type Application interface {
	DefaultFlagSet() *flag.FlagSet
	Run(context.Context) (interface{}, error)
	RunWithFlagSet(context.Context, *flag.FlagSet) (interface{}, error)
}
