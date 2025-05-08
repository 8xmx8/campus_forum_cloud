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
	BaseUrl    string       // API 基础地址
	HTTPClient *http.Client // HTTP 客户端
	Model      string       // 模型名称
}

func NewOllamaClient(baseURL, model string) *Client {
	return &Client{
		Model:   model,
		BaseUrl: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10000 * time.Second,
		},
	}
}

func (c *Client) CommentMessage(message string) (*Response, error) {
	request := &Request{
		Model: c.Model,
		Messages: []*RoleContent{
			{
				Role:    "system",
				Content: "用户向你输入一段话，请你鉴别是否涉及不良言论，如果是，请给出1到5之间的一个数字和理由，严格按照下面的格式输出，强调一下只有 1 表示用户正常聊天，其余的 2，3，4，5 这个几个数字表示用户不良言论。比如 1 | 用户正常交流。2 | 用户不良言论。3 | 用户对亲人进行辱骂。3 | 用户涉及敏感话题。4 | 用户大量不良言论。 5 | 用户试图对人身进行攻击",
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
		return nil, err
	}
	var response *Response
	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}
	return response, nil
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
	Messages []*RoleContent `json:"messages"`
	Stream   bool           `json:"stream"`
}
type RoleContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Response struct {
	Model              string       `json:"model"`
	CreatedAt          time.Time    `json:"created_at"`
	Messages           *RoleContent `json:"message"`
	DoneReason         string       `json:"done_reason"`
	Done               bool         `json:"done"`
	TotalDuration      int64        `json:"total_duration"`
	LoadDuration       int64        `json:"load_duration"`
	PromptEvalCount    int          `json:"prompt_eval_count"`
	PromptEvalDuration int64        `json:"prompt_eval_duration"`
	EvalCount          int          `json:"eval_count"`
	EvalDuration       int64        `json:"eval_duration"`
}
