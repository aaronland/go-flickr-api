package auth

// https://www.flickr.com/services/api/auth.oauth.html

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type OAuth1RequestToken struct {
	RequestToken
	OAuthToken       string `json:"oath_token"`
	OAuthTokenSecret string `json:"oauth_token_secret"`
}

func (t *OAuth1RequestToken) Token() string {
	return t.OAuthToken
}

func (t *OAuth1RequestToken) Secret() string {
	return t.OAuthTokenSecret
}

type OAuth1AuthorizationToken struct {
	AuthorizationToken
	OAuthToken    string `json:"oath_token"`
	OAuthVerifier string `json:"oath_verifier"`
}

func (t *OAuth1AuthorizationToken) Token() string {
	return t.OAuthToken
}

func (t *OAuth1AuthorizationToken) Verifier() string {
	return t.OAuthVerifier
}

type OAuth1AccessToken struct {
	AccessToken
	OAuthToken       string `json:"oauth_token"`
	OAuthTokenSecret string `json:"oauth_token_secret"`
}

func (t *OAuth1AccessToken) Token() string {
	return t.OAuthToken
}

func (t *OAuth1AccessToken) Secret() string {
	return t.OAuthTokenSecret
}

func UnmarshalOAuth1RequestToken(str_q string) (RequestToken, error) {

	q, err := url.ParseQuery(str_q)

	if err != nil {
		return nil, err
	}

	required := []string{
		"oauth_token",
		"oauth_token_secret",
	}

	_, err = ensureQueryParameters(q, required...)

	if err != nil {
		return nil, err
	}

	tok := &OAuth1RequestToken{
		OAuthToken:       q.Get("oauth_token"),
		OAuthTokenSecret: q.Get("oauth_token_secret"),
	}

	return tok, nil
}

func UnmarshalOAuth1AuthorizationToken(str_q string) (AuthorizationToken, error) {

	q, err := url.ParseQuery(str_q)

	if err != nil {
		return nil, err
	}

	required := []string{
		"oauth_token",
		"oauth_verifier",
	}

	_, err = ensureQueryParameters(q, required...)

	if err != nil {
		return nil, err
	}

	tok := &OAuth1AuthorizationToken{
		OAuthToken:    q.Get("oauth_token"),
		OAuthVerifier: q.Get("oauth_verifier"),
	}

	return tok, nil
}

func UnmarshalOAuth1AccessToken(str_q string) (AccessToken, error) {

	q, err := url.ParseQuery(str_q)

	if err != nil {
		return nil, err
	}

	required := []string{
		"oauth_token",
		"oauth_token_secret",
	}

	_, err = ensureQueryParameters(q, required...)

	if err != nil {
		return nil, err
	}

	tok := &OAuth1AccessToken{
		OAuthToken:       q.Get("oauth_token"),
		OAuthTokenSecret: q.Get("oauth_token_secret"),
	}

	return tok, nil
}

func ensureQueryParameters(query url.Values, keys ...string) (bool, error) {

	for _, k := range keys {

		v := query.Get(k)

		if v == "" {
			return false, fmt.Errorf("Missing '%s' parameter", k)
		}
	}

	return true, nil
}

/*

First, you must create a base string from your request. The base string is constructed by concatenating the HTTP verb,
the request URL, and all request parameters sorted by name, using lexicograhpical byte value ordering, separated by an '&'.

*/

func GenerateOAuth1SigningBaseString(http_method string, endpoint *url.URL, args *url.Values) string {

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

func GenerateOAuth1Signature(key string, base string) string {

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return ret
}
