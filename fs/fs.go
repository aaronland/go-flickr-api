package fs

import (
	"context"
	"fmt"
	"io"
	io_fs "io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aaronland/go-flickr-api/client"
	"github.com/tidwall/gjson"
)

type apiFS struct {
	io_fs.FS
	http_client *http.Client
	client      client.Client
}

func New(ctx context.Context, cl client.Client) io_fs.FS {

	http_cl := &http.Client{}

	fs := &apiFS{
		http_client: http_cl,
		client:      cl,
	}

	return fs
}

func (f *apiFS) Open(name string) (io_fs.File, error) {

	args := &url.Values{}
	args.Set("method", "flickr.photos.getInfo")
	args.Set("photo_id", name)

	ctx := context.Background()
	r, err := f.client.ExecuteMethod(ctx, args)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	id_rsp := gjson.GetBytes(body, "photo.id")

	if !id_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.id")
	}

	secret_rsp := gjson.GetBytes(body, "photo.secret")

	if !secret_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.secret")
	}

	originalsecret_rsp := gjson.GetBytes(body, "photo.originalsecret")

	if !originalsecret_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.originalsecret")
	}

	originalformat_rsp := gjson.GetBytes(body, "photo.originalformat")

	if !originalformat_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.originalformat")
	}

	server_rsp := gjson.GetBytes(body, "photo.server")

	if !server_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.server")
	}

	farm_rsp := gjson.GetBytes(body, "photo.farm")

	if !farm_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.farm")
	}

	lastmod_rsp := gjson.GetBytes(body, "photo.dates.lastupdate")

	if !lastmod_rsp.Exists() {
		return nil, fmt.Errorf("Missing photo.dates.lastupdate")
	}

	id := id_rsp.Int()
	// secret := secret_rsp.String()
	originalsecret := originalsecret_rsp.String()
	originalformat := originalformat_rsp.String()
	server := server_rsp.String()
	lastmod := lastmod_rsp.Int()

	url := fmt.Sprintf("https://live.staticflickr.com/%s/%d_%s_o.%s", server, id, originalsecret, originalformat)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new request, %w", err)
	}

	rsp, err := f.http_client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute request, %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		defer rsp.Body.Close()
		return nil, fmt.Errorf("%d %s", rsp.StatusCode, rsp.Status)
	}

	str_len := rsp.Header.Get("Content-Length")
	int_len, _ := strconv.ParseInt(str_len, 10, 64)

	t := time.Unix(lastmod, 0)

	// To do: Derive file permissions from Flickr permissions
	// "visibility":{"ispublic":0,"isfriend":0,"isfamily":0}

	fl := &File{
		name:           filepath.Base(url),
		content:        rsp.Body,
		content_length: int_len,
		modTime:        t,
	}

	return fl, nil
}

func (f apiFS) ReadFile(name string) ([]byte, error) {
	r, err := f.Open(name)

	if err != nil {
		return nil, err
	}

	defer r.Close()
	return io.ReadAll(r)
}

func (f *apiFS) ReadDir(name string) ([]io_fs.DirEntry, error) {

	ctx := context.Background()

	args, err := url.ParseQuery(name)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse query, %w", err)
	}

	entries := []io_fs.DirEntry{}

	cb := func(ctx context.Context, r io.ReadSeekCloser, err error) error {

		defer r.Close()

		if err != nil {
			return err
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read API response body, %w", err)
		}

		fmt.Println(string(body))
		return nil
	}

	err = client.ExecuteMethodPaginatedWithClient(ctx, f.client, &args, cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	return entries, nil
}

func (f *apiFS) Sub(path string) (io_fs.FS, error) {
	return nil, fmt.Errorf("Not supported")
}
