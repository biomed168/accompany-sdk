//go:build js && wasm

package wasm_wrapper

import (
	"accompany-sdk/pkg/utils"
	"accompany-sdk/sdk"
	"accompany-sdk/wasm/event_listener"
	"syscall/js"
)

const COMMONEVENTFUNC = "commonEventFunc"

type WrapperCommon struct {
	commonFunc *js.Value
}

// NewWrapperCommon 处理来自 JavaScript 的事件或函数调用，并存储一个 JavaScript 函数（或值）的引用
func NewWrapperCommon() *WrapperCommon {
	return &WrapperCommon{}
}

// CommonEventFunc 接收一个回调函数
func (w *WrapperCommon) CommonEventFunc(_ js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		w.commonFunc = &args[len(args)-1]
		return js.ValueOf(true)
	} else {
		return js.ValueOf(false)
	}
}

type WrapperInit struct {
	*WrapperCommon
}

func NewWrapperInit(wrapperCommon *WrapperCommon) *WrapperInit {
	return &WrapperInit{WrapperCommon: wrapperCommon}
}

func (w *WrapperInit) InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewConnCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return js.ValueOf(event_listener.NewCaller(sdk.InitSDK, callback, &args).SyncCall())
}

func (w *WrapperInit) Login(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(sdk.Login, callback, &args).AsyncCallWithCallback()
}
