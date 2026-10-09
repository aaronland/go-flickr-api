package sanitize

import (
	"net/http"
	"strconv"

	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

// PostString returns a sanitised string extracted from the POST form
// body of the given request for the specified key.  If the key is missing
// an empty string is returned.
func PostString(req *http.Request, param string) (string, error) {

	raw_value := req.PostFormValue(param)
	return wof_sanitize.SanitizeString(raw_value, sn_opts)
}

// PostInt returns an int extracted from the POST form body of the
// given request for the specified key.  The value is first sanitised
// and then parsed as int.  If the key is missing or the value is
// empty, zero is returned.
func PostInt(req *http.Request, param string) (int, error) {

	str_value, err := PostString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.Atoi(str_value)
}

// PostInt64 returns an int64 extracted from the POST form body of the
// given request for the specified key.  The value is first sanitised
// and then parsed as int64.  If the key is missing or the value is
// empty, zero is returned.
func PostInt64(req *http.Request, param string) (int64, error) {

	str_value, err := PostString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseInt(str_value, 10, 64)
}

// PostFloat64 returns a float64 extracted from the POST form body of
// the given request for the specified key.  The value is first
// sanitised and then parsed as float64.  If the key is missing or the
// value is empty, zero is returned.
func PostFloat64(req *http.Request, param string) (float64, error) {

	str_value, err := PostString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseFloat(str_value, 64)
}

// PostBool returns a bool extracted from the POST form body of the
// given request for the specified key.  The value is first sanitised
// and then parsed as bool.  If the key is missing or the value is
// empty, false is returned.
func PostBool(req *http.Request, param string) (bool, error) {

	str_value, err := PostString(req, param)

	if err != nil {
		return false, err
	}

	if str_value == "" {
		return false, nil
	}

	return strconv.ParseBool(str_value)
}
