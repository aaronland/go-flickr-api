package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/whosonfirst/go-ioutil"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

const API string = "https://www.flickr.com/services"
const AUTH_REQUEST = "oauth/request_token"
const AUTH_AUTHORIZE = "oauth/authorize"
const AUTH_TOKEN = "oauth/request_token"

const REST string = "rest"

type HTTPClient struct {
	Client
	http_client        *http.Client
	api_endpoint       *url.URL
	consumer_key       string
	consumer_secret    string
	oauth_token        string
	oauth_token_secret string
}

func NewHTTPClient(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	key := q.Get("consumer_key")
	secret := q.Get("consumer_secret")

	if key == "" {
		return nil, fmt.Errorf("Missing ?consumer_key parameter")
	}

	if secret == "" {
		return nil, fmt.Errorf("Missing ?consumer_secret parameter")
	}

	http_client := &http.Client{}

	cl := &HTTPClient{
		http_client:     http_client,
		consumer_key:    key,
		consumer_secret: secret,
	}

	return cl, nil
}

func (cl *HTTPClient) GetRequestToken(ctx context.Context, cb_url *url.URL) (*auth.RequestToken, error) {

	endpoint, err := url.Parse(API)

	if err != nil {
		return nil, err
	}

	endpoint.Path = filepath.Join(endpoint.Path, AUTH_REQUEST)

	http_method := "GET"

	args := &url.Values{}
	args.Set("oauth_callback", cb_url.String())

	args, err = cl.signArgs(http_method, endpoint, args)

	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	fh, err := cl.call(ctx, req)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	rsp_body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return auth.UnmarshalRequestToken(string(rsp_body))
}

func (cl *HTTPClient) AuthorizationURL(ctx context.Context, req *auth.RequestToken) (*url.URL, error) {

	q := url.Values{}
	q.Set("oauth_token", req.Token)

	u, err := url.Parse(API)

	if err != nil {
		return nil, err
	}

	u.Path = AUTH_AUTHORIZE
	u.RawQuery = q.Encode()

	return u, nil
}

func (cl *HTTPClient) GetAccessToken(ctx context.Context, auth_token *auth.AuthorizationToken) (*auth.AccessToken, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (cl *HTTPClient) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(API)

	if err != nil {
		return nil, err
	}

	endpoint.Path = filepath.Join(endpoint.Path, REST)

	http_method := "GET"

	args, err = cl.prepareArgs(http_method, endpoint, args)

	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	return cl.call(ctx, req)
}

func (cl *HTTPClient) call(ctx context.Context, req *http.Request) (io.ReadSeekCloser, error) {

	req = req.WithContext(ctx)

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}

func (cl *HTTPClient) prepareArgs(http_method string, endpoint *url.URL, args *url.Values) (*url.Values, error) {

	args.Set("nojsoncallback", "1")
	args.Set("format", "json")

	return cl.signArgs(http_method, endpoint, args)
}

func (cl *HTTPClient) signArgs(http_method string, endpoint *url.URL, args *url.Values) (*url.Values, error) {

	now := time.Now()
	ts := now.Unix()

	str_ts := strconv.FormatInt(ts, 10)

	nonce := auth.GenerateNonce()

	args.Set("oauth_version", "1.0")
	args.Set("oauth_signature_method", "HMAC-SHA1")

	args.Set("oauth_nonce", nonce)
	args.Set("oauth_timestamp", str_ts)
	args.Set("oauth_consumer_key", cl.consumer_key)

	if cl.oauth_token != "" {
		args.Set("oauth_token", cl.oauth_token)
	}

	// args.Set("api_key", cl.consumer_key)

	sig := cl.getSignature(http_method, endpoint, args, cl.oauth_token_secret)
	args.Set("oauth_signature", sig)

	return args, nil
}

func (cl *HTTPClient) getSignature(http_method string, endpoint *url.URL, args *url.Values, token_secret string) string {

	key := fmt.Sprintf("%s&%s", url.QueryEscape(cl.consumer_secret), url.QueryEscape(token_secret))
	base_string := auth.GenerateSigningBaseString(http_method, endpoint, args)

	return auth.GenerateSignature(key, base_string)
}
