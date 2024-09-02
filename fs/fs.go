package fs

import (
	"context"
	"fmt"
	"io"
	io_fs "io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aaronland/go-flickr-api/client"
	"github.com/tidwall/gjson"
)

var re_photo = regexp.MustCompile(`^(?:\d+|(?:.*?\#)?\/\d+\/\d+_\w+_[a-z]\.\w+)$`)
var re_url = regexp.MustCompile(`\#?(\/\d+\/\d+_\w+_[a-z]\.\w+)$`)

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

	ctx := context.Background()

	logger := slog.Default()
	logger = logger.With("name", name)

	logger.Debug("Open file")

	if !re_photo.MatchString(name) {

		logger.Debug("File does not match photo ID or URL, assuming SPR entry")

		fl := &apiFile{
			name:           name,
			content_length: -1,
			modTime:        time.Now(),
			is_spr:         true,
		}

		return fl, nil
	}

	var path string

	if re_url.MatchString(name) {

		m := re_url.FindStringSubmatch(name)
		path = m[1]
		logger.Debug("Derive relative path", "rel path", path)
	} else {

		args := &url.Values{}
		args.Set("method", "flickr.photos.getInfo")
		args.Set("photo_id", name)

		logger.Debug("Get photo info", "query", args.Encode())

		r, err := f.client.ExecuteMethod(ctx, args)

		if err != nil {
			return nil, fmt.Errorf("Failed to execute API method, %w", err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			return nil, fmt.Errorf("Failed to read API response body, %w", err)
		}

		// logger.Debug("API response", "body", string(body))

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
		// lastmod := lastmod_rsp.Int()

		path = fmt.Sprintf("/%s/%d_%s_o.%s", server, id, originalsecret, originalformat)
	}

	u, err := url.Parse("https://live.staticflickr.com")

	if err != nil {
		return nil, fmt.Errorf("Failed to parse base URL (which is weird), %w", err)
	}

	u.Path = path
	url := u.String()

	logger.Debug("Fetch URL", "url", url)

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

	// last-modified: Sat, 31 Aug 2024 20:07:47 GMT
	lastmod := rsp.Header.Get("Last-Modified")

	t, err := time.Parse(time.RFC1123, lastmod)

	if err != nil {
		logger.Error("Failed to parse lastmod time, default to now", "lastmod", lastmod, "error", err)
		t = time.Now()
	}

	// To do: Derive file permissions from Flickr permissions
	// "visibility":{"ispublic":0,"isfriend":0,"isfamily":0}

	fl := &apiFile{
		name:           u.Path,
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

	logger := slog.Default()
	logger = logger.With("name", name)

	ctx := context.Background()

	args, err := url.ParseQuery(name)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse query, %w", err)
	}

	extras := make([]string, 0)

	ensure_extras := []string{
		"url_o",
		"lastupdate",
	}

	if args.Has("extras") {
		extras = strings.Split(args.Get("extras"), ",")
		args.Del("extras")
	}

	for _, v := range ensure_extras {

		if !slices.Contains(extras, v) {
			extras = append(extras, v)
		}
	}

	args.Set("extras", strings.Join(extras, ","))

	logger.Debug("Read dir")

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

		logger.Debug(string(body))

		rsp := gjson.GetBytes(body, "*.photo")

		if !rsp.Exists() {
			return fmt.Errorf("Failed to derive photos from response")
		}

		// {\"id\":\"53961664838\",\"secret\":\"a11f30f8e0\",\"server\":\"65535\",\"farm\":66,\"title\":\"113028507763072851\",\"isprimary\":\"0\",\"ispublic\":0,\"isfriend\":0,\"isfamily\":0}

		for _, ph := range rsp.Array() {

			url_rsp := ph.Get("url_o")

			if !url_rsp.Exists() {
				logger.Warn("Response is missing url_o extra, skipping")
				continue
			}

			url_o, err := url.Parse(url_rsp.String())

			if err != nil {
				return fmt.Errorf("Failed to parse url_o value (%s), %w", url_rsp.String(), err)
			}

			lastmod_rsp := ph.Get("lastupdate")
			lastmod := time.Unix(lastmod_rsp.Int(), 0)

			fi := &apiFileInfo{
				name:    fmt.Sprintf("#%s", url_o.Path),
				size:    -1,
				is_spr:  false,
				modTime: lastmod,
			}

			ent := &apiDirEntry{
				info: fi,
			}

			logger.Debug("Add entry", "path", fi.name)
			entries = append(entries, ent)
		}

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
