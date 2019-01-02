package mfh

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var defhc http.Client

func init() {
	defhc = http.Client{}
	defhc.Timeout = time.Minute * 10
}

// DefaultClient default client
func DefaultClient() (c http.Client) {
	return defhc
}

// HTTPCall - simple call http
func HTTPCall(method string, url string, headers map[string]string,
	cookies []*http.Cookie, timeout time.Duration, body io.Reader) (data []byte,
	statusCode int,
	status string,
	headersRes map[string][]string,
	cookiesRes map[string]string,
	err error) {

	headersRes = make(map[string][]string)
	cookiesRes = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	if cookies != nil {
		for _, v := range cookies {
			req.AddCookie(v)
		}
	}

	r2 := req.WithContext(ctx)

	res, err := defhc.Do(r2)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	for k, v := range res.Header {
		if len(v) > 0 {
			headersRes[k] = v
		}
	}

	for _, v := range res.Cookies() {
		cookiesRes[v.Name] = v.Value
	}

	statusCode = res.StatusCode
	status = res.Status

	return
}
