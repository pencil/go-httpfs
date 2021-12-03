package httpfs_test

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pencil/go-httpfs/httpfs"
)

type httpClientMock struct {
	do func(req *http.Request) (*http.Response, error)
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return h.do(req)
}

func TestRead(t *testing.T) {
	content := "Hello World"
	f := httpfs.New(
		&httpClientMock{
			do: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path != "/base/foo" {
					t.Errorf("expected path /base/foo, got %s", req.URL.Path)
				}
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(content)),
				}, nil
			},
		},
		&url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/base/",
		},
	)
	res, err := f.Open("foo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()

	data, err := io.ReadAll(res)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("expected %s, got %s", content, data)
	}
}

func TestStat(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		header          http.Header
		contentLength   int64
		expectedName    string
		expectedSize    int64
		expectedModTime time.Time
	}{
		{
			name: "with all expected headers present",
			path: "yeet/foo.html",
			header: http.Header{
				"Last-Modified": []string{"Wed, 01 Jan 2000 00:00:00 GMT"},
			},
			contentLength:   123,
			expectedName:    "foo.html",
			expectedSize:    123,
			expectedModTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "without expected headers present",
			path:            "x/y/z.js",
			header:          http.Header{},
			contentLength:   -1,
			expectedName:    "z.js",
			expectedSize:    -1,
			expectedModTime: time.Time{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := "Hello World"
			f := httpfs.New(
				&httpClientMock{
					do: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode:    200,
							Header:        test.header,
							ContentLength: test.contentLength,
							Body:          io.NopCloser(strings.NewReader(content)),
						}, nil
					},
				},
				&url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/",
				},
			)
			res, err := f.Open(test.path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Close()

			stat, err := res.Stat()
			if err != nil {
				t.Fatal(err)
			}

			if stat.Name() != test.expectedName {
				t.Errorf("expected name %s, got %s", test.expectedName, stat.Name())
			}

			if stat.Size() != test.expectedSize {
				t.Errorf("expected size %d, got %d", test.expectedSize, stat.Size())
			}

			if !stat.ModTime().Equal(test.expectedModTime) {
				t.Errorf("expected mod time %s, got %s", test.expectedModTime, stat.ModTime())
			}

			if stat.IsDir() {
				t.Error("expected not a directory")
			}

			if stat.Mode() != 0 {
				t.Errorf("expected mode 0, got %d", stat.Mode())
			}

			if stat.Sys() != nil {
				t.Error("expected no sys")
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	f := httpfs.New(
		&httpClientMock{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
		},
		&url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/",
		},
	)
	_, err := f.Open("foo")
	if err != fs.ErrNotExist {
		t.Errorf("expected ErrNotExist, got %v", err)
	}
}

func TestServerError(t *testing.T) {
	f := httpfs.New(
		&httpClientMock{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     "500 Internal Server Error",
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
		},
		&url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/",
		},
	)
	_, err := f.Open("foo")
	if err == nil {
		t.Error("expected error")
	} else if err.Error() != "HTTP error: 500 Internal Server Error" {
		t.Errorf("expected 500 Internal Server Error, got %v", err)
	}
}

func TestHTTPClientError(t *testing.T) {
	expectedErr := errors.New("something went wrong")
	f := httpfs.New(
		&httpClientMock{
			do: func(req *http.Request) (*http.Response, error) {
				return nil, expectedErr
			},
		},
		&url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/",
		},
	)
	_, err := f.Open("foo")
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
