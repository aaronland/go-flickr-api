package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

/*

First, you must create a base string from your request. The base string is constructed by concatenating the HTTP verb,
the request URL, and all request parameters sorted by name, using lexicograhpical byte value ordering, separated by an '&'.

*/

func GenerateSigningBaseString(http_method string, endpoint *url.URL, args *url.Values) string {

	endpoint_url := endpoint.String()
	request_url := url.QueryEscape(endpoint_url)

	enc_args := args.Encode()
	flickr_encoded := strings.Replace(enc_args, "+", "%20", -1)

	query := url.QueryEscape(flickr_encoded)

	ret := fmt.Sprintf("%s&%s&%s", http_method, request_url, query)
	return ret
}

/*

Use the base string as the text and the key is the concatenated values of the Consumer Secret and Token Secret, separated by an '&'.

*/

func GenerateSignature(key string, base string) string {

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return ret
}
