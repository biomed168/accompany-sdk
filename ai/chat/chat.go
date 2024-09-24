package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/mylxsw/aidea-server/pkg/misc"
	"strings"

	"github.com/mylxsw/aidea-server/config"
	"github.com/mylxsw/go-utils/array"
)

var (
	ErrContextExceedLimit = errors.New("上下文长度超过最大限制")
	ErrContentFilter      = errors.New("请求或响应内容包含敏感词")
)

type Message struct {
	Role              string              `json:"role"`
	Content           string              `json:"content"`
	MultipartContents []*MultipartContent `json:"multipart_content,omitempty"`
}

func (m Message) UploadedFile() *FileURL {
	mc := array.Filter(m.MultipartContents, func(part *MultipartContent, _ int) bool {
		return part.Type == "file" && part.FileURL != nil && part.FileURL.URL != ""
	})

	if len(mc) == 0 {
		return nil
	}

	return mc[0].FileURL
}

type MultipartContent struct {
	// Type 对于 OpenAI 来说， type 可选值为 image_url/text
	Type     string    `json:"type"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
	Text     string    `json:"text,omitempty"`
	FileURL  *FileURL  `json:"file_url,omitempty"`
}

type FileURL struct {
	// URL is the URL of the file
	URL string `json:"url,omitempty"`
	// Name is the name of the file
	Name string `json:"name,omitempty"`
}

type ImageURL struct {
	// URL Either a URL of the image or the base64 encoded image data.
	URL string `json:"url,omitempty"`
	// Detail Specifies the detail level of the image
	// Three options, low, high, or auto, you have control over how the model processes the image and generates its textual understanding.
	// By default, the model will use the auto setting which will look at the image input size and decide if it should use the low or high setting
	//
	// - `low` will disable the “high res” model. The model will receive a low-res 512px x 512px version of the image,
	//   and represent the image with a budget of 65 tokens. This allows the API to return faster responses and consume
	//   fewer input tokens for use cases that do not require high detail.
	//
	// - `high` will enable “high res” mode, which first allows the model to see the low res image and
	//   then creates detailed crops of input images as 512px squares based on the input image size.
	//   Each of the detailed crops uses twice the token budget (65 tokens) for a total of 129 tokens.
	Detail string `json:"detail,omitempty"`
}

type Messages []Message

// Purification Filter out unsupported message types (file document upload), currently only support text and image_url.
func (ms Messages) Purification() Messages {
	return array.Map(ms, func(item Message, _ int) Message {
		if len(item.MultipartContents) <= 0 {
			return item
		}

		mul := make([]*MultipartContent, 0)
		for _, content := range item.MultipartContents {
			if array.In(content.Type, []string{"image_url", "text"}) {
				mul = append(mul, content)
				continue
			}
		}

		item.MultipartContents = mul
		return item
	})
}

func (ms Messages) ToLogEntry() Messages {
	ret := make(Messages, len(ms))
	for i, msg := range ms {
		mm := Message{
			Role:    msg.Role,
			Content: misc.SubString(msg.Content, 20),
		}

		if msg.MultipartContents != nil {
			mm.MultipartContents = make([]*MultipartContent, len(msg.MultipartContents))
			for j, part := range msg.MultipartContents {
				mm.MultipartContents[j] = &MultipartContent{
					Type: part.Type,
					Text: misc.SubString(part.Text, 20),
				}

				if part.ImageURL != nil {
					mm.MultipartContents[j].ImageURL = &ImageURL{
						URL:    misc.SubString(part.ImageURL.URL, 20),
						Detail: part.ImageURL.Detail,
					}
				}
			}
		}

		ret[i] = mm
	}

	return ret
}

func (ms Messages) HasImage() bool {
	for _, msg := range ms {
		for _, part := range msg.MultipartContents {
			if part.ImageURL != nil && part.ImageURL.URL != "" {
				return true
			}
		}
	}

	return false
}

// Fix 模型上下文预处理：
// 1. 强制上下文为 user/assistant 轮流出现
// 2. 第一个普通消息必须是用户消息
// 3. 最后一条消息必须是用户消息
func (ms Messages) Fix() Messages {
	msgs := ms
	// 如果最后一条消息不是用户消息，则补充一条用户消息
	last := msgs[len(msgs)-1]
	if last.Role != "user" {
		last = Message{
			Role:    "user",
			Content: "继续",
		}
		msgs = append(msgs, last)
	}

	// 过滤掉 system 消息，因为 system 消息需要在每次对话中保留，不受上下文长度限制
	systemMsgs := array.Filter(msgs, func(m Message, _ int) bool { return m.Role == "system" })
	if len(systemMsgs) > 0 {
		msgs = array.Filter(msgs, func(m Message, _ int) bool { return m.Role != "system" })
	}

	finalMessages := make([]Message, 0)
	var lastRole string

	for _, m := range array.Reverse(msgs) {
		if m.Role == lastRole {
			continue
		}

		lastRole = m.Role
		finalMessages = append(finalMessages, m)
	}

	if len(finalMessages)%2 == 0 {
		finalMessages = finalMessages[:len(finalMessages)-1]
	}

	return append(systemMsgs, array.Reverse(finalMessages)...)
}

// Request represents a request structure for chat completion API.
type Request struct {
	Stream    bool     `json:"stream,omitempty"`
	Model     string   `json:"model"`
	Messages  Messages `json:"messages"`
	MaxTokens int      `json:"max_tokens,omitempty"`
	N         int      `json:"n,omitempty"` // 复用作为 room_id

	// 业务定制字段
	RoomID    int64 `json:"-"`
	WebSocket bool  `json:"-"`

	// TempModel 用户可以指定临时模型来进行当前对话，实现临时切换模型的功能
	TempModel string `json:"temp_model,omitempty"`
}

func (req Request) Clone() Request {
	return Request{
		Stream:    req.Stream,
		Model:     req.Model,
		Messages:  array.Map(req.Messages, func(item Message, _ int) Message { return item }),
		MaxTokens: req.MaxTokens,
		N:         req.N,
		RoomID:    req.RoomID,
		WebSocket: req.WebSocket,
		TempModel: req.TempModel,
	}
}

func (req Request) Purification() Request {
	req.Messages = req.Messages.Purification()
	return req
}

func (req Request) assembleMessage() string {
	var msgs []string
	for _, msg := range req.Messages {
		msgs = append(msgs, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
	}

	return strings.Join(msgs, "\n\n")
}

func (req Request) Init() Request {
	// 去掉模型名称前缀
	modelSegs := strings.Split(req.Model, ":")
	if len(modelSegs) > 1 {
		modelSegs = modelSegs[1:]
	}

	req.Model = strings.Join(modelSegs, ":")

	// 获取 room id
	// 这里复用了参数 N
	req.RoomID = int64(req.N)
	if req.N != 0 {
		req.N = 0
	}

	// 过滤掉内容为空的 message
	req.Messages = array.Filter(req.Messages, func(item Message, _ int) bool { return strings.TrimSpace(item.Content) != "" })

	return req
}

// Fix 修复请求内容，注意：上下文长度修复后，最终的上下文数量不包含 system 消息和用户最后一条消息
func (req Request) Fix(chat Chat, maxContextLength int64, maxTokenCount int) (*Request, int64, error) {
	// 自动缩减上下文长度至满足模型要求的最大长度，尽可能避免出现超过模型上下文长度的问题
	systemMessages := array.Filter(req.Messages, func(item Message, _ int) bool { return item.Role == "system" })
	systemMessageLen, _ := MessageTokenCount(systemMessages, req.Model)

	// 模型允许的 Tokens 数量和请求参数指定的 Tokens 数量，取最小值
	modelTokenLimit := chat.MaxContextLength(req.Model) - systemMessageLen
	if modelTokenLimit < maxTokenCount {
		maxTokenCount = modelTokenLimit
	}

	messages, inputTokens, err := ReduceMessageContext(
		ReduceMessageContextUpToContextWindow(
			array.Filter(req.Messages, func(item Message, _ int) bool { return item.Role != "system" }),
			int(maxContextLength),
		),
		req.Model,
		maxTokenCount,
	)
	if err != nil {
		return nil, 0, errors.New("超过模型最大允许的上下文长度限制，请尝试“新对话”或缩短输入内容长度")
	}

	req.Messages = array.Map(append(systemMessages, messages...), func(item Message, _ int) Message {
		if len(item.MultipartContents) > 0 {
			item.MultipartContents = array.Map(item.MultipartContents, func(part *MultipartContent, _ int) *MultipartContent {
				if part.ImageURL != nil && part.ImageURL.URL != "" && part.ImageURL.Detail == "" {
					part.ImageURL.Detail = "low"
				}

				return part
			})
		}
		return item
	})

	return &req, int64(inputTokens), nil
}

func (req Request) ResolveCalFeeModel(conf *config.Config) string {
	return req.Model
}

type Response struct {
	Error        string `json:"error,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	Text         string `json:"text,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

type Chat interface {
	// Chat 以请求-响应的方式进行对话
	Chat(ctx context.Context, req Request) (*Response, error)
	// ChatStream 以流的方式进行对话
	ChatStream(ctx context.Context, req Request) (<-chan Response, error)
	// MaxContextLength 获取模型的最大上下文长度
	MaxContextLength(model string) int
}

type AI struct {
	OpenAI *OpenAIChat
	Baidu  *BaiduAIChat
}
