package ccontext

import (
	"accompany-sdk/ai_struct"
	"accompany-sdk/sdk_callback"
	"accompany-sdk/sdk_struct"
	"context"

	"github.com/openimsdk/tools/mcontext"
)

const (
	Callback = "callback"
)

type GlobalConfig struct {
	UserID string
	Token  string

	sdk_struct.SDKConfig
}

type ContextInfo interface {
	UserID() string
	Token() string
	PlatformID() int32
	DataDir() string
	LogLevel() uint32
	OperationID() string
	IsExternalExtensions() bool

	OpenAIConfig() *ai_struct.OpenAiConfig
}

func Info(ctx context.Context) ContextInfo {
	conf := ctx.Value(GlobalConfigKey{}).(*GlobalConfig)
	return &info{
		conf: conf,
		ctx:  ctx,
	}
}

func WithInfo(ctx context.Context, conf *GlobalConfig) context.Context {
	return context.WithValue(ctx, GlobalConfigKey{}, conf)
}

func WithOperationID(ctx context.Context, operationID string) context.Context {
	return mcontext.SetOperationID(ctx, operationID)
}
func WithSendMessageCallback(ctx context.Context, callback sdk_callback.SendMsgCallBack) context.Context {
	return context.WithValue(ctx, Callback, callback)
}

func WithApiErrCode(ctx context.Context, cb ApiErrCodeCallback) context.Context {
	return context.WithValue(ctx, apiErrCode{}, cb)
}

func GetApiErrCodeCallback(ctx context.Context) ApiErrCodeCallback {
	fn, _ := ctx.Value(apiErrCode{}).(ApiErrCodeCallback)
	if fn == nil {
		return &emptyApiErrCodeCallback{}
	}
	return fn
}

type GlobalConfigKey struct{}

type info struct {
	conf *GlobalConfig
	ctx  context.Context
}

func (i *info) UserID() string {
	return i.conf.UserID
}

func (i *info) Token() string {
	return i.conf.Token
}

func (i *info) PlatformID() int32 {
	return i.conf.PlatformID
}

func (i *info) DataDir() string {
	return i.conf.DataDir
}

func (i *info) LogLevel() uint32 {
	return i.conf.LogLevel
}

func (i *info) OperationID() string {
	return mcontext.GetOperationID(i.ctx)
}

func (i *info) IsExternalExtensions() bool {
	return i.conf.IsExternalExtensions
}

// OpenAIConfig 返回openai的配置
func (i *info) OpenAIConfig() *ai_struct.OpenAiConfig {
	return &i.conf.AiConfig.OpenAiConfig
}

type apiErrCode struct{}

type ApiErrCodeCallback interface {
	OnError(ctx context.Context, err error)
}

type emptyApiErrCodeCallback struct{}

func (e *emptyApiErrCodeCallback) OnError(ctx context.Context, err error) {}
