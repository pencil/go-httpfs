# go-httpfs
An [`fs.FS`](https://pkg.go.dev/io/fs#FS) implementation that reads files from an HTTP backend.

## Usage

```go
package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"

	"github.com/pencil/go-httpfs/httpfs"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:3000/data/")
	var store fs.FS
	store = httpfs.New(http.DefaultClient, baseURL)

	file, err := store.Open("test.txt") // requests http://localhost:3000/data/test.txt
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", data)
}
```
