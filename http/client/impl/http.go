package impl

import (
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/http/client/request"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	HttpSchema                 = "http://"
)

// HttpClient http客户端,支持长连接
type HttpClient struct {
	Timeout    time.Duration
	Transport  *http.Transport
	clientPool *http.Client
}

// NewHttpClient 默认http客户端
func NewHttpClient() HttpClient {
	return HttpClient{
		Timeout: time.Duration(3) * time.Second,
		Transport: &http.Transport{ // 配置连接池
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        DefaultMaxIdleConns,
			MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(90) * time.Second,
		},
	}
}

// NewRequest 实例化request
func (client HttpClient) NewRequest(param request.Param) *http.Request {
	if !strings.HasPrefix(param.Url, HttpSchema) {
		param.Url = HttpSchema + param.Url
	}

	req, err := http.NewRequest(param.Method, param.Url, param.Body)
	if err != nil {
		return nil
	}

	return req
}

// Execute 执行请求
func (client HttpClient) Execute(request request.IRequest) (string, error) {
	if request == nil {
		return "", errors.New("request is nil")
	}

	req, ok := request.(*http.Request)

	if !ok {
		return "", errors.New("request not *http.request")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	resp, err := client.clientPool.Do(req)
	if err != nil {
		return "", err
	}

	respStr, err := helper.StringifyResponse(resp)
	if err != nil {
		return "", err
	}
	return respStr, nil

}
