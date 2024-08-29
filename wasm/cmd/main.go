package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"syscall/js"

	"accompany-sdk/wasm/wasm_wrapper"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MAIN", "panic info is:", r, debug.Stack())
		}
	}()
	fmt.Println("runtime env", runtime.GOARCH, runtime.GOOS)
	registerFunc()
	<-make(chan bool)
}

func registerFunc() {
	// 注册全局回调函数，用于给js通知内容
	globalFuc := wasm_wrapper.NewWrapperCommon()
	js.Global().Set(wasm_wrapper.COMMONEVENTFUNC, js.FuncOf(globalFuc.CommonEventFunc))

	// 注册一个 init Wrapper
	wrapperInit := wasm_wrapper.NewWrapperInit(globalFuc)
	js.Global().Set("initSDK", js.FuncOf(wrapperInit.InitSDK))
	js.Global().Set("login", js.FuncOf(wrapperInit.Login))
}
