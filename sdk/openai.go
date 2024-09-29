package sdk

import (
	"accompany-sdk/sdk_callback"
)

// AskOpenAi 简单问询
func AskOpenAi(callback sdk_callback.Base, operationID string, prompt string, question string, maxTokenCount int) {
	call(callback, operationID, UserForSDK.OpenAi().QuickAsk, prompt, question, maxTokenCount)
}
