package resty

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	client *resty.Client
	once   sync.Once
)

func initRty() {
	client = resty.New().
		SetTransport(&http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}).
		// 设置连接超时
		SetTimeout(10 * time.Second).
		// SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetCookieJar(nil)
	client.OnAfterResponse(RestyLogMiddleware)
}

func GetRty() *resty.Client {
	once.Do(initRty)
	return client
}

func RestyLogMiddleware(_ *resty.Client, resp *resty.Response) error {
	if resp.IsError() {
		method := resp.Request.Method
		url := resp.Request.URL
		zap.L().Error("请求出现错误", zap.String("method", method),
			zap.String("url", url), zap.Int64("time_spent(ms)", resp.Time().Milliseconds()))
	}
	return nil
}
