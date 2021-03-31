package client

import (
	_ "bytes"
	"context"
	"crypto/hmac"
	_ "crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	_ "sort"
	"strconv"
	"strings"
	"time"
)

const API string = "https://www.flickr.com/services"
const AUTH_REQUEST = "oauth/request_token"
const AUTH_AUTHORIZE = "oauth/authorize"
const AUTH_TOKEN = "oauth/request_token"

const REST string = "rest"

type HTTPClient struct {
	Client
	http_client     *http.Client
	api_endpoint    *url.URL
	consumer_key    string
	consumer_secret string
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

	api, err := url.Parse(API)

	if err != nil {
		return nil, err
	}

	api.Path = REST

	http_client := &http.Client{}

	cl := &HTTPClient{
		http_client:     http_client,
		api_endpoint:    api,
		consumer_key:    key,
		consumer_secret: secret,
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

	http_method := "POST"

	args, err := cl.prepareArgs(http_method, args)

	if err != nil {
		return nil, err
	}

	str_args := args.Encode()
	body := strings.NewReader(str_args)

	url := cl.api_endpoint.String()

	log.Println(url, str_args)

	req, err := http.NewRequest(http_method, url, body)

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

func (cl *HTTPClient) prepareArgs(http_method string, args *url.Values) (*url.Values, error) {

	args.Set("nojsoncallback", "1")
	args.Set("format", "json")

	return cl.signArgs(http_method, args)
}

func (cl *HTTPClient) signArgs(http_method string, args *url.Values) (*url.Values, error) {

	now := time.Now()
	ts := now.Unix()

	str_ts := strconv.FormatInt(ts, 10)

	nonce := generateNonce()
	sig := "FIX ME"

	args.Set("oauth_version", "1.0")
	args.Set("oauth_timestamp", str_ts)
	args.Set("oauth_consumer_key", cl.consumer_key)
	args.Set("oauth_signature_method", "HMAC-SHA1")
	args.Set("oauth_nonce", nonce)

	args.Set("oauth_signature", sig)

	return args, nil
}

// The following are all cribbed from
// https://github.com/masci/flickr/blob/v2/client.go

// Get the base string to compose the signature
func (cl *HTTPClient) getSigningBaseString(http_method string, args *url.Values) string {

	endpoint_url := cl.api_endpoint.String()
	request_url := url.QueryEscape(endpoint_url)

	flickr_encoded := strings.Replace(args.Encode(), "+", "%20", -1)
	query := url.QueryEscape(flickr_encoded)

	ret := fmt.Sprintf("%s&%s&%s", http_method, request_url, query)
	return ret
}

// Compute the signature of a signed request
func (cl *HTTPClient) getSignature(http_method string, args *url.Values, token_secret string) string {

	key := fmt.Sprintf("%s&%s", url.QueryEscape(cl.consumer_secret), url.QueryEscape(token_secret))

	base_string := cl.getSigningBaseString(http_method, args)

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(base_string))

	ret := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return ret
}

/*
// Sign API requests. This method differs from the signing process needed for
// OAuth authenticated requests.
func (cl *HTTPClient) getApiSignature(token_secret string) string {

	var buf bytes.Buffer
	buf.WriteString(token_secret)

	keys := make([]string, 0, len(c.Args))
	for k := range c.Args {
		keys = append(keys, k)
	}
	// args needs to be in alphabetical order
	sort.Strings(keys)

	for _, k := range keys {
		arg := c.Args[k][0]
		buf.WriteString(k)
		buf.WriteString(arg)
	}

	base := buf.String()

	data := []byte(base)
	return fmt.Sprintf("%x", md5.Sum(data))
}
*/

// Generate a random string of 8 chars, needed for OAuth signature
func generateNonce() string {

	rand.Seed(time.Now().UTC().UnixNano())

	// For convenience, use a set of chars we don't need to url-escape
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
