package httpclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)


type Response struct {
	StatusCode int
	Status     string
	Body       string
	Duration   time.Duration
	Err        string
}

func Send(method, url, body string, headers map[string]string) Response {
	start := time.Now()

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
		if headers == nil {
			headers = make(map[string]string)
		}
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/json"
		}
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return Response{Err: err.Error()}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Response{Err: err.Error(), Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{Err: err.Error(), Duration: time.Since(start)}
	}

	return Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       strings.TrimSpace(string(raw)),
		Duration:   time.Since(start),
	}
}
