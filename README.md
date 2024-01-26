# `net/http.RoundTripper` for `acw_sc__v2`

## Usage

```bash
go get github.com/WangYihang/acw-sc-v2
```

## Example

```go
package main

import (
	"fmt"
	"io"
	"net/http"

	acw_sc_v2 "github.com/WangYihang/acw-sc-v2"
)

func main() {
	client := &http.Client{
		Transport: acw_sc_v2.NewTransport(), // Use acw_sc_v2 as Transport
	}
	for i := 0; i < 8; i++ {
		resp, err := client.Get("https://www.example.com/")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(body)
	}
}
```

## References

* NodeJS version for Server API ([acw-sc-v2.js](https://github.com/WangYihang/acw-sc-v2.js))
* Python version for Client Code ([acw-sc-v2-py](https://github.com/WangYihang/acw-sc-v2-py))