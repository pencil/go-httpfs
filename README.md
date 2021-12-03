# go-httpfs
An `fs.FS` implementation that reads files from an HTTP backend.

## Usage

```go
baseURL, _ := url.Parse("http://localhost:3000/data")
var store fs.FS
store = httpfs.New(http.DefaultClient, baseURL)

file, err := store.Open("test.txt") // requests http://localhost:3000/data/test.txt
if err != nil {
  // ...
}

data, err := io.ReadAll(file)
if err != nil {
  // ...
}

fmt.Println(data)
```
