package httpfs

import (
	"errors"
	"io/fs"
	"net/http"
	"net/url"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpFS struct {
	httpClient httpClient
	baseURL    *url.URL
}

var _ fs.FS = (*httpFS)(nil) // Ensure HTTPFS implements FS.

func New(httpClient httpClient, baseURL *url.URL) *httpFS {
	return &httpFS{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (f *httpFS) Open(name string) (fs.File, error) {
	req := &http.Request{
		Method: "GET",
		URL:    f.baseURL.ResolveReference(&url.URL{Path: name}),
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 400 {
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil, fs.ErrNotExist
		}
		return nil, errors.New("HTTP error: " + resp.Status)
	}

	return &httpFile{
		req:  req,
		resp: resp,
	}, nil
}
