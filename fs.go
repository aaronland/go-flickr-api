package api

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aaronland/go-flickr-api/client"
	"github.com/tidwall/gjson"
)

type FS struct {
	fs.FS
	http_client *http.Client
	client      client.Client
}

func NewFS(ctx context.Context, cl client.Client) *FS {

	http_cl := &http.Client{}

	fs := &FS{
		http_client: http_cl,
		client:      cl,
	}

	return fs
}

func (f *FS) Open(name string) (fs.File, error) {

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

	// "visibility":{"ispublic":0,"isfriend":0,"isfamily":0}

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
		return nil, err
	}

	rsp, err := f.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		defer rsp.Body.Close()
		return nil, fmt.Errorf("%d %s", rsp.StatusCode, rsp.Status)
	}

	str_len := rsp.Header.Get("Content-Length")
	int_len, _ := strconv.ParseInt(str_len, 10, 64)

	t := time.Unix(lastmod, 0)

	fl := File{
		name:           filepath.Base(url),
		content:        rsp.Body,
		content_length: int_len,
		modTime:        t,
	}

	return fl, nil
}

type File struct {
	name           string
	perm           os.FileMode
	content        io.ReadCloser
	content_length int64
	modTime        time.Time
	closed         bool
}

func (f *File) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	fi := fileInfo{
		name:    f.name,
		size:    f.content_length,
		modTime: f.modTime,
		mode:    f.perm,
	}
	return &fi, nil
}

func (f *File) Read(b []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	return f.content.Read(b)
}

func (f *File) Close() error {
	if f.closed {
		return fs.ErrClosed
	}
	err := f.content.Close()

	if err != nil {
		return err
	}

	f.closed = true
	return nil
}

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    fs.FileMode
}

// base name of the file
func (fi *fileInfo) Name() string {
	return fi.name
}

// length in bytes for regular files; system-dependent for others
func (fi *fileInfo) Size() int64 {
	return fi.size
}

// file mode bits
func (fi *fileInfo) Mode() fs.FileMode {
	return fi.mode
}

// modification time
func (fi *fileInfo) ModTime() time.Time {
	return fi.modTime
}

// abbreviation for Mode().IsDir()
func (fi *fileInfo) IsDir() bool {
	return false
}

// underlying data source (can return nil)
func (fi *fileInfo) Sys() interface{} {
	return nil
}
