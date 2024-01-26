package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	acw_sc_v2 "github.com/WangYihang/acw-sc-v2"
)

func createRequest() *http.Request {
	request, _ := http.NewRequest("GET", "http://www.beianx.cn/search/ctycdn.com", nil)
	var headers = map[string]string{
		`Accept-Encoding`:           `identity`,
		`Accept`:                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8application/signed-exchange;v=b3;q=0.7`,
		`Accept-Language`:           `en`,
		`Cache-Control`:             `max-age=0`,
		`Connection`:                `keep-alive`,
		`Referer`:                   `https://www.beianx.cn/search`,
		`Sec-Fetch-Dest`:            `document`,
		`Sec-Fetch-Mode`:            `navigate`,
		`Sec-Fetch-Site`:            `same-origin`,
		`Upgrade-Insecure-Requests`: `1`,
		`User-Agent`:                `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0`,
		`sec-ch-ua`:                 `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`,
		`sec-ch-ua-mobile`:          `?0`,
		`sec-ch-ua-platform`:        `"Windows"`,
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}

func main() {
	client := &http.Client{
		Transport: acw_sc_v2.NewTransport(), // Use acw_sc_v2 as Transport
	}
	request := createRequest()
	for i := 0; i < 8; i++ {
		// Send request
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		// Read response body
		rawBody, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		// Close response body
		response.Body.Close()

		// Print response body
		fmt.Println(string(rawBody))

		// Sleep for 1 second
		time.Sleep(time.Duration(1) * time.Second)
	}
}
