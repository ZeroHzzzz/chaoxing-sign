package resty

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

var (
	client *resty.Client
	once   sync.Once
)

func initRty() {
	client = resty.New().
		// SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetCookieJar(nil)
}

func GetRty() *resty.Client {
	once.Do(initRty)
	return client
}

func HttpSendPost(url string, req map[string]any, headers map[string]string, resp any) (*resty.Response, error) {
	client := GetRty()

	r, err := client.R().
		SetHeaders(headers).
		SetBody(req).
		SetResult(&resp).
		Post(url)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func HttpSendGet(url string, headers map[string]string, query map[string]string, body map[string]any, resp any) (*resty.Response, error) {
	client := GetRty()

	req := client.R()

	if headers != nil {
		req.SetHeaders(headers)
	}

	if query != nil {
		req.SetQueryParams(query)
	}

	if body != nil {
		req.SetBody(body)
	}

	r, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	return r, nil
}
