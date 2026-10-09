package sanitize

import (
	"net/http"
)

func RequestStringMulti(req *http.Request, param string) ([]string, error) {

	switch req.Method {

	case "POST":
		return PostStringMulti(req, param)
	default:
		return GetStringMulti(req, param)
	}

}

func RequestIntMulti(req *http.Request, param string) ([]int, error) {

	switch req.Method {

	case "POST":
		return PostIntMulti(req, param)
	default:
		return GetIntMulti(req, param)
	}

}

func RequestInt64Multi(req *http.Request, param string) ([]int64, error) {

	switch req.Method {

	case "POST":
		return PostInt64Multi(req, param)
	default:
		return GetInt64Multi(req, param)
	}

}

func RequestFloat64Multi(req *http.Request, param string) ([]float64, error) {

	switch req.Method {

	case "POST":
		return PostFloat64Multi(req, param)
	default:
		return GetFloat64Multi(req, param)
	}

}

func RequestBoolMulti(req *http.Request, param string) ([]bool, error) {

	switch req.Method {

	case "POST":
		return PostBoolMulti(req, param)
	default:
		return GetBoolMulti(req, param)
	}

}
