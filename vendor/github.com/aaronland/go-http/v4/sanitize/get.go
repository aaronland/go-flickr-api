package sanitize

import (
	"net/http"
	"strconv"

	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

// GetString returns a sanitised string extracted from the URL query
// parameters of the given request for the specified key.  If the key
// is missing an empty string is returned.
func GetString(req *http.Request, param string) (string, error) {

	q := req.URL.Query()
	raw_value := q.Get(param)
	return wof_sanitize.SanitizeString(raw_value, sn_opts)
}

// GetInt returns an int extracted from the URL query parameters of the
// given request for the specified key.  The value is first sanitised
// and then parsed as an int.  If the key is missing or the value is
// empty, zero is returned.
func GetInt(req *http.Request, param string) (int, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.Atoi(str_value)
}

// GetInt64 returns an int64 extracted from the URL query parameters of
// the given request for the specified key.  The value is first
// sanitised and then parsed as an int64.  If the key is missing or the
// value is empty, zero is returned.
func GetInt64(req *http.Request, param string) (int64, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseInt(str_value, 10, 64)
}

// GetFloat64 returns a float64 extracted from the URL query parameters
// of the given request for the specified key.  The value is first
// sanitised and then parsed as a float64.  If the key is missing or the
// value is empty, zero is returned.
func GetFloat64(req *http.Request, param string) (float64, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseFloat(str_value, 64)
}

// GetBool returns a bool extracted from the URL query parameters of
// the given request for the specified key.  The value is first
// sanitised and then parsed as a bool.  If the key is missing or the
// value is empty, false is returned.
func GetBool(req *http.Request, param string) (bool, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return false, err
	}

	if str_value == "" {
		return false, nil
	}

	return strconv.ParseBool(str_value)
}
