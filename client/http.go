package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const API string = "https://www.flickr.com/services"
const AUTH_REQUEST = "oauth/request_token"
const AUTH_AUTHORIZE = "oauth/authorize"
const AUTH_TOKEN = "oauth/request_token"

const REST string = "rest"

type HTTPClient struct {
	Client
	http_client  *http.Client
	api_endpoint *url.URL
	consumer_key string
}

func NewHTTPClient(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	key := q.Get("consumer_key")

	if key == "" {
		return nil, fmt.Errorf("Missing ?consumer_key parameter")
	}

	api, err := url.Parse(API)

	if err != nil {
		return nil, err
	}

	api.Path = REST

	http_client := &http.Client{}

	cl := &HTTPClient{
		http_client:  http_client,
		api_endpoint: api,
		consumer_key: key,
	}

	return cl, nil
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

func (cl *HTTPClient) GetRequestToken(ctx context.Context, cb_url *url.URL) (*auth.RequestToken, error) {

	args := &url.Values{}
	args.Set("oauth_callback", cb_url.String())

	fh, err := cl.ExecuteMethod(ctx, args)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return auth.UnmarshalRequestToken(string(body))
}

func (cl *HTTPClient) GetAccessToken(ctx context.Context, auth_token *auth.AuthorizationToken) (*auth.AccessToken, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (cl *HTTPClient) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	args, err := cl.prepareArgs(args)

	if err != nil {
		return nil, err
	}

	body := strings.NewReader(args.Encode())

	url := cl.api_endpoint.String()

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp)
}

func (cl *HTTPClient) prepareArgs(args *url.Values) (*url.Values, error) {

	args.Set("nojsoncallback", "1")
	args.Set("format", "json")

	return auth.SignArgs(cl.consumer_key, args)
}
