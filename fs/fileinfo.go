package fs

import (
	io_fs "io/fs"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    io_fs.FileMode
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
func (fi *fileInfo) Mode() io_fs.FileMode {
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
