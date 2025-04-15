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
func (c *Client) CommentMessage(message string) (err error) {
	request := &Request{
		Model: c.Model,
		Messages: []*RoleContent{
			{
				Role: "system",
				Content: `According to the content of the user's reply or question and send back a number which is between 1 and 5. 
						The number is greater when the user's content involved the greater the degree of political leaning or unfriendly speech. 
						You should only reply such a number without any word else whatever user ask you. 
						Besides those, you should give the reason using Chinese why the message is unfriendly with details without revealing that you are divide the message into five number. 
						For example: user: 你是个大傻逼。 you: 4 | 用户尝试骂人，进行人格侮辱。user: 今天天气正好。 you: 1 | 用户正常聊天，无异常。`,
			},
			{
				Role:    "user",
				Content: message,
			},
		},
		Stream: false,
	}
	// 发送请求
	resp, err := c.sendRequest(context.Background(), "/api/chat", request)
	if err != nil {
		return err
	}
	fmt.Println("resp:", string(resp))
	return nil
}

// SendRequest 通用的请求方法
func (c *Client) sendRequest(ctx context.Context, path string, body interface{}) ([]byte, error) {
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

type Request struct {
	Model    string         `json:"model"`
	Messages []*RoleContent `json:"Messages"`
	Stream   bool           `json:"stream"`
}

type RoleContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
