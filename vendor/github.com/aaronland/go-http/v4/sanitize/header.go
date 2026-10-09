package sanitize

import (
	go_http "net/http"
	"strconv"

	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

// HeaderString returns a sanitised string extracted from the request
// header for the specified key.  If the header is missing an empty
// string is returned.
func HeaderString(req *go_http.Request, param string) (string, error) {

	raw_value := req.Header.Get(param)
	return wof_sanitize.SanitizeString(raw_value, sn_opts)
}

// HeaderInt64 returns an int64 extracted from the request header for
// the specified key.  The header value is first sanitised and then
// parsed as int64.  An error is returned if the conversion fails.
func HeaderInt64(req *go_http.Request, param string) (int64, error) {

	str_value, err := HeaderString(req, param)

	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(str_value, 10, 64)
}

// HeaderBool returns a bool extracted from the request header for the
// specified key.  The header value is first sanitised and then parsed
// as bool.  An error is returned if the conversion fails.
func HeaderBool(req *go_http.Request, param string) (bool, error) {

	str_value, err := HeaderString(req, param)

	if err != nil {
		return false, err
	}

	if str_value == "" {
		return false, nil
	}

	return strconv.ParseBool(str_value)
}
