package httpfs

import (
	"io/fs"
	"net/http"
)

type httpFile struct {
	req  *http.Request
	resp *http.Response
}

var _ fs.File = (*httpFile)(nil) // Ensure httpFile implements File.

func (f *httpFile) Read(p []byte) (int, error) {
	return f.resp.Body.Read(p)
}

func (f *httpFile) Close() error {
	return f.resp.Body.Close()
}

func (f *httpFile) Stat() (fs.FileInfo, error) {
	return &httpFileInfo{
		req:  f.req,
		resp: f.resp,
	}, nil
}
