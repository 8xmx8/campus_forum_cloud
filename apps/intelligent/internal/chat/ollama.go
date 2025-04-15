package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var ErrTimeout = errors.New("request timeout")                // 请求超时
var ServerAddrIsEmpty = errors.New("ollama address is empty") // ollama 地址为空

type Client struct {
	BaseUrl    string            // API 基础地址
	HTTPClient *http.Client      // HTTP 客户端
	Headers    map[string]string // 公共请求头
	Model      string            // 模型名称
}

func NewOllamaClient(baseURL, model string) *Client {
	return &Client{
		Model:   model,
		BaseUrl: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
func (c *Client) CommentMessage() {
}

// SendRequest 通用的请求方法
func (c *Client) SendRequest(ctx context.Context, path string, body interface{}) ([]byte, error) {
	if c.BaseUrl == "" {
		return nil, ServerAddrIsEmpty
	}
	api := c.BaseUrl + path
	// 序列化请求体（如果有）
	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
		fmt.Println("请求:", string(jsonBody))
	} else {
		bodyReader = bytes.NewReader(nil)
	}
	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// 检查错误是否是 context deadline exceeded
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		// 检查错误是否是 *url.Error
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			// 检查底层错误是否是超时
			var netErr net.Error
			if errors.As(urlErr.Err, &netErr) && netErr.Timeout() {
				return nil, ErrTimeout
			}
		} else {
			return nil, err
		}
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp.Body == nil {
		return nil, errors.New("response body is nil")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
