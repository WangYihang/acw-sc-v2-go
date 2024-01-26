package acw_sc_v2

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewTransport() http.RoundTripper {
	t := &antiScrapeTransport{original: http.DefaultTransport}
	t.original.(*http.Transport).DisableKeepAlives = true
	return t
}

type antiScrapeTransport struct {
	acw_sc__v2 string
	original   http.RoundTripper
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

func ReadAllBody(response *http.Response) []byte {
	var reader io.ReadCloser
	var err error

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var rawBodyBuffer bytes.Buffer
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			slog.Error("error occured while reading response body", slog.String("error", err.Error()))
		}
		buffer := make([]byte, 1024)
		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				rawBodyBuffer.Write(buffer[0:n])
			}
			if err != nil {
				slog.Error("error occured while reading response body", slog.String("error", err.Error()))
				break
			}
		}
		reader.Close()
		return rawBodyBuffer.Bytes()
	default:
		reader = response.Body
		rawBody, _ := io.ReadAll(reader)
		response.Body.Close()
		return rawBody
	}
}

func (t *antiScrapeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check if there is a valid cookie
	if t.acw_sc__v2 != "" {
		slog.Info("using existing cookie", slog.String("cookie", t.acw_sc__v2))
		req.AddCookie(&http.Cookie{
			Name:    "acw_sc__v2",
			Value:   t.acw_sc__v2,
			MaxAge:  3600,
			Path:    "/",
			Expires: time.Now().Add(3600 * time.Second),
		})
	}
	// Send the original request
	resp, err := t.original.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// Read the response body to check if anti scrape is triggered
	rawBody := ReadAllBody(resp)
	if isAntiScrapeTriggered(rawBody) {
		// Anti scrape is triggered
		slog.Info("anti-scrape detected")
		// Generate new cookie
		cookieValue, err := crackTheJSCodeAndGetCookie(rawBody)
		if err != nil {
			slog.Error("error occured while cracking the js code: %v", err)
		}
		// Set new cookie to the current request
		req.AddCookie(&http.Cookie{Name: "acw_sc__v2", Value: cookieValue})
		t.acw_sc__v2 = cookieValue
		// Send the original request again
		slog.Info("resending the original request")
		return t.original.RoundTrip(req)
	} else {
		// Anti scrape is not triggered
		slog.Info("anti-scrape not detected")
		resp.Body = io.NopCloser(bytes.NewReader(rawBody))
		return resp, nil
	}
}
