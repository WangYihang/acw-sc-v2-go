package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	acw_sc_v2 "github.com/WangYihang/acw-sc-v2"
)

func main() {
	client := &http.Client{
		Transport: acw_sc_v2.NewTransport(), // Use acw_sc_v2 as Transport
	}
	for i := 0; i < 8; i++ {
		// Send request
		resp, err := client.Get("https://www.example.com/")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		// Print response body
		fmt.Println(body)
		// Sleep for 1 second
		time.Sleep(time.Second)
	}
}
