// Package sanitize provides helper functions for extracting and sanitising
// parameters from HTTP requests.  The package supports retrieving values
// from URL query strings, POST form bodies and request headers, and
// converting them to common Go types (string, int, int64, float64, bool).
// All string values are sanitised using the options defined in the
// package variable `sn_opts`.
package sanitize

import (
	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

var sn_opts *wof_sanitize.Options

func init() {
	sn_opts = wof_sanitize.DefaultOptions()
}
