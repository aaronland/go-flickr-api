package fs

import (
	"io"
	io_fs "io/fs"
	"os"
	"time"
)

type File struct {
	name           string
	perm           os.FileMode
	content        io.ReadCloser
	content_length int64
	modTime        time.Time
	closed         bool
}

func (f *File) Stat() (io_fs.FileInfo, error) {

	if f.closed {
		return nil, io_fs.ErrClosed
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
		return 0, io_fs.ErrClosed
	}

	return f.content.Read(b)
}

func (f *File) Close() error {

	if f.closed {
		return io_fs.ErrClosed
	}

	err := f.content.Close()

	if err != nil {
		return err
	}

	f.closed = true
	return nil
}
