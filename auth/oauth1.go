package auth

import (
	"net/url"
	"strconv"
	"time"
)

const OAUTH_VERSION string = "1.0"

func SignArgs(key string, args *url.Values) (*url.Values, error) {

	now := time.Now()
	ts := now.Unix()

	str_ts := strconv.FormatInt(ts, 10)

	nonce := "FIX ME"
	sig := "FIX ME"

	args.Set("oauth_version", OAUTH_VERSION)
	args.Set("oauth_timestamp", str_ts)
	args.Set("oauth_consumer_key", key)
	args.Set("oauth_signature_method", "HMAC-SHA1")
	args.Set("oauth_nonce", nonce)

	args.Set("oauth_signature", sig)

	return args, nil
}
