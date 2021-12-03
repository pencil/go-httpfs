# go-httpfs
An `fs.FS` implementation that reads files from an HTTP backend.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"

	"github.com/pencil/go-httpfs/httpfs"
)

func main() {
	println("Hello, World!")
	baseURL, _ := url.Parse("http://localhost:3000/data/")
	var store fs.FS
	store = httpfs.New(http.DefaultClient, baseURL)

	file, err := store.Open("test.txt") // requests http://localhost:3000/data/test.txt
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
```
