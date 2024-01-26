package acw_sc_v2

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func NewTransport() http.RoundTripper {
	return &antiScrapeTransport{original: http.DefaultTransport}
}

type antiScrapeTransport struct {
	original http.RoundTripper
}

func isAntiScrapeTriggered(rawBody []byte) bool {
	return bytes.Contains(rawBody, []byte("acw_sc__v2"))
}

func crackTheJSCodeAndGetCookie(rawBody []byte) (string, error) {
	data := []string{string(rawBody)}
	formData := url.Values{
		"data": data,
	}

	endpoint := "https://acw-sc-v2.authu.online/"
	slog.Info("cracking the js code", slog.String("endpoint", endpoint), slog.String("body", string(rawBody[0:32])))
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookie, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	slog.Info("cookie generated", slog.String("cookie", string(cookie)))

	return string(cookie), nil
}

func readAllBody(response *http.Response) []byte {
	var reader io.ReadCloser
	var err error
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			slog.Error("error occured while reading response body", slog.String("error", err.Error()))
		}
	default:
		slog.Info("default")
		reader = response.Body
	}
	defer reader.Close()

	buffer := make([]byte, 1024)
	rawBodyBuffer := bytes.Buffer{}
	for {
		for i := range buffer {
			buffer[i] = 0
		}
		n, err := reader.Read(buffer)
		rawBodyBuffer.Write(buffer[0:n])
		if err != nil {
			slog.Error("error occured while reading response body", slog.String("error", err.Error()))
			break
		}
	}
	return rawBodyBuffer.Bytes()
}

func (t *antiScrapeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.original.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	rawBody := readAllBody(resp)
	if isAntiScrapeTriggered(rawBody) {
		slog.Info("anti-scrape detected")
		cookieValue, err := crackTheJSCodeAndGetCookie(rawBody)
		if err != nil {
			slog.Error("error occured while cracking the js code: %v", err)
		}
		cookieJar, _ := cookiejar.New(nil)
		cookieURL, _ := url.Parse(req.URL.String())
		cookieJar.SetCookies(cookieURL, []*http.Cookie{
			{Name: "acw_sc__v2", Value: cookieValue},
		})
		client := &http.Client{Transport: t.original, Jar: cookieJar}
		slog.Info("resending the original request")
		return client.Do(req)
	}
	return resp, nil
}
