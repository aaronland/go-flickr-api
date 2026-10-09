package sanitize

import (
	"net/http"
	"strconv"

	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

// GetStringMulti returns a slice of sanitised strings extracted from the
// URL query parameters of the given request for the specified key.
// If the key is not present an empty slice is returned without error.
func GetStringMulti(req *http.Request, param string) ([]string, error) {

	q := req.URL.Query()
	raw_values, ok := q[param]

	if !ok {
		return make([]string, 0), nil
	}

	values := make([]string, len(raw_values))

	for i, v := range raw_values {

		sv, err := wof_sanitize.SanitizeString(v, sn_opts)

		if err != nil {
			return nil, err
		}

		values[i] = sv
	}

	return values, nil
}

// GetIntMulti returns a slice of integers extracted from the URL query
// parameters of the given request for the specified key.  The function
// first sanitises the raw string values and then converts each to an
// `int`.  If the key is not present an empty slice is returned.
func GetIntMulti(req *http.Request, param string) ([]int, error) {

	string_values, err := GetStringMulti(req, param)

	if err != nil {
		return nil, err
	}

	int_values := make([]int, len(string_values))

	for i, v := range string_values {

		sv, err := strconv.Atoi(v)

		if err != nil {
			return nil, err
		}

		int_values[i] = sv
	}

	return int_values, nil
}

// GetInt64Multi returns a slice of int64 values extracted from the URL
// query parameters of the given request for the specified key.  The raw
// values are sanitised before being parsed as int64.
func GetInt64Multi(req *http.Request, param string) ([]int64, error) {

	string_values, err := GetStringMulti(req, param)

	if err != nil {
		return nil, err
	}

	int_values := make([]int64, len(string_values))

	for i, v := range string_values {

		sv, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return nil, err
		}

		int_values[i] = sv
	}

	return int_values, nil
}

// GetFloat64Multi returns a slice of float64 values extracted from the
// URL query parameters of the given request for the specified key.  The
// raw values are sanitised before being parsed as float64.
func GetFloat64Multi(req *http.Request, param string) ([]float64, error) {

	string_values, err := GetStringMulti(req, param)

	if err != nil {
		return nil, err
	}

	fl_values := make([]float64, len(string_values))

	for i, v := range string_values {

		sv, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return nil, err
		}

		fl_values[i] = sv
	}

	return fl_values, nil
}

// GetBoolMulti returns a slice of bool values extracted from the URL query
// parameters of the given request for the specified key.  The raw values
// are sanitised before being parsed as bool.
func GetBoolMulti(req *http.Request, param string) ([]bool, error) {

	string_values, err := GetStringMulti(req, param)

	if err != nil {
		return nil, err
	}

	bool_values := make([]bool, len(string_values))

	for i, v := range string_values {

		sv, err := strconv.ParseBool(v)

		if err != nil {
			return nil, err
		}

		bool_values[i] = sv
	}

	return bool_values, nil
}
