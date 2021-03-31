package auth

// https://www.flickr.com/services/api/auth.oauth.html

import (
	"fmt"
	"net/url"
)

type RequestToken struct {
	Token  string `json:"oath_token"`
	Secret string `json:"oauth_token_secret"`
}

type AuthorizationToken struct {
	Token    string `json:"oath_token"`
	Verifier string `json:"oath_verifier"`
}

type AccessToken struct {
	Token    string `json:"oath_token"`
	Secret   string `json:"oauth_token_secret"`
	NSID     string `json:"user_nsid"`
	Username string `json:"username"`
}

func UnmarshalRequestToken(str_q string) (*RequestToken, error) {

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

	tok := &RequestToken{
		Token:  q.Get("oauth_token"),
		Secret: q.Get("oauth_token_secret"),
	}

	return tok, nil
}

func UnmarshalAuthorizationToken(str_q string) (*AuthorizationToken, error) {

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

	tok := &AuthorizationToken{
		Token:    q.Get("oauth_token"),
		Verifier: q.Get("oauth_verifier"),
	}

	return tok, nil
}

func UnmarshalAccessToken(str_q string) (*AccessToken, error) {

	q, err := url.ParseQuery(str_q)

	if err != nil {
		return nil, err
	}

	required := []string{
		"oauth_token",
		"oauth_token_secret",
		"user_nsid",
		"username",
	}

	_, err = ensureQueryParameters(q, required...)

	if err != nil {
		return nil, err
	}

	tok := &AccessToken{
		Token:    q.Get("oauth_token"),
		Secret:   q.Get("oauth_token_secret"),
		NSID:     q.Get("user_nsid"),
		Username: q.Get("username"),
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
