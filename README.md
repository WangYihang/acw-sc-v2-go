# `net/http.RoundTripper` for `acw_sc__v2`

## Usage

```bash
go get github.com/WangYihang/acw-sc-v2
```

## Example

```go
package main

import (
	"net/http"

	acw_sc_v2 "github.com/WangYihang/acw-sc-v2"
)

func main() {
	client := &http.Client{
		Transport: acw_sc_v2.NewTransport(), // Use acw_sc_v2 as Transport
	}
	resp, err := client.Get("https://www.example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
```