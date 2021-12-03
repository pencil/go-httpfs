package httpfs

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"time"
)

type httpFileInfo struct {
	req  *http.Request
	resp *http.Response
}

var _ fs.FileInfo = (*httpFileInfo)(nil) // Ensure httpFileInfo implements FileInfo.

func (f *httpFileInfo) Name() string {
	return filepath.Base(f.req.URL.Path)
}

func (f *httpFileInfo) Size() int64 {
	return f.resp.ContentLength
}

func (f *httpFileInfo) Mode() fs.FileMode {
	return 0
}

func (f *httpFileInfo) ModTime() time.Time {
	lm := f.resp.Header.Get("Last-Modified")
	if lm == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC1123, lm)
	return t
}

func (f *httpFileInfo) IsDir() bool {
	return false
}

func (f *httpFileInfo) Sys() interface{} {
	return nil
}
