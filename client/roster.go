package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type ClientInitializeFunc func(context.Context, string) (Client, error)

var clients roster.Roster

func ensureClientRoster() error {

	if clients == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		clients = r
	}

	return nil
}

func RegisterClient(ctx context.Context, scheme string, f ClientInitializeFunc) error {

	err := ensureClientRoster()

	if err != nil {
		return err
	}

	return clients.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureClientRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range clients.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewClient(ctx context.Context, uri string) (Client, error) {

	// To account for things that might be gocloud.dev/runtimevar-encoded
	// in a file using editors that automatically add newlines (thanks, Emacs)

	uri = strings.TrimSpace(uri)

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := clients.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ClientInitializeFunc)
	return f(ctx, uri)
}
