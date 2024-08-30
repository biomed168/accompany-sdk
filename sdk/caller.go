// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"accompany-sdk/pkg/ccontext"
	"accompany-sdk/pkg/sdkerrs"
	"accompany-sdk/sdk_callback"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
)

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func setNumeric(in interface{}, out interface{}) {
	inValue := reflect.ValueOf(in)
	outValue := reflect.ValueOf(out)
	outElem := outValue.Elem()
	outType := outElem.Type()
	inType := inValue.Type()
	if outType.AssignableTo(inType) {
		outElem.Set(inValue)
		return
	}
	inKind := inValue.Kind()
	outKind := outElem.Kind()
	switch inKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(inValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Int()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Uint()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(inValue.Uint())
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Uint()))
		}
	case reflect.Float32, reflect.Float64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Float()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Float()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(inValue.Float())
		}
	}
}

// call_ 根据传入的函数指针 `fn` 和参数 `args` 动态调用该函数。
// 参数：
//   - operationID: 用于标识调用操作的唯一ID。
//   - fn: 需要被调用的函数指针。
//   - args: 变长参数，用于传递给被调用函数的参数列表。
// 返回值：
//   - res: 被调用函数的返回结果。
//   - err: 调用过程中产生的错误，如果有的话。

// 函数执行流程：
// 1. 使用 `recover` 机制捕获调用过程中可能出现的 panic 并返回相应的错误信息。
// 2. 通过反射获取函数指针及函数名称，并进行初步的有效性检查。
// 3. 调用 `CheckResourceLoad` 检查资源是否已经加载，未加载则返回错误。
// 4. 创建一个上下文 `ctx`，用于记录和追踪该操作的日志。
// 5. 检查 `fn` 参数是否是一个可调用的函数，如果不是，则返回错误。
// 6. 验证 `fn` 的参数数量与提供的 `args` 数量是否匹配。
// 7. 根据 `fn` 的参数类型和 `args` 的实际类型，准备好实际的调用参数：
//   - 如果 `args` 的类型不匹配 `fn` 的参数类型，尝试进行必要的类型转换。
//   - 如果是字符串类型，尝试解析 JSON 并解码为相应的 Go 类型。
//
// 8. 使用反射 `Call` 方法调用 `fn`，并获取返回值。
// 9. 检查返回的最后一个值是否实现了 `error` 接口，如果是且不为 `nil`，则记录日志并返回错误。
// 10. 处理函数返回值的 `nil` 情况（如空的 map 或 slice），并将它们初始化为对应的类型。
// 11. 根据函数返回的结果数量，返回单个结果或者多个结果的集合，并记录日志。
func call_(operationID string, fn any, args ...any) (res any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("call panic: %+v", r)
		}
	}()
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	if operationID == "" {
		return nil, sdkerrs.ErrArgs.WrapMsg("call function operationID is empty")
	}
	if err := CheckResourceLoad(UserForSDK, funcName); err != nil {
		return nil, sdkerrs.ErrResourceLoad.WrapMsg("not load resource")
	}
	ctx := ccontext.WithOperationID(UserForSDK.BaseCtx(), operationID)
	log.ZInfo(ctx, "call function", "in sdk args", args)

	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("call function fn is not function, is %T", fn))
	}
	fnt := fnv.Type()
	nin := fnt.NumIn()
	if len(args)+1 != nin {
		return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("go code error: fn in args num is not match"))
	}
	t := time.Now()
	log.ZInfo(ctx, "input req", "function name", funcName, "args", args)
	ins := make([]reflect.Value, 0, nin)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		inFnField := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == inFnField.String() || inFnField.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			var ptr int
			for inFnField.Kind() == reflect.Ptr {
				inFnField = inFnField.Elem()
				ptr++
			}
			switch inFnField.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
				v := reflect.New(inFnField)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("go call json.Unmarshal error: %s", err))
				}
				if ptr == 0 {
					v = v.Elem()
				} else if ptr != 1 {
					for i := ptr - 1; i > 0; i-- {
						temp := reflect.New(v.Type())
						temp.Elem().Set(v)
						v = temp
					}
				}
				ins = append(ins, v)
				continue
			}
		}
		//if isNumeric(arg.Kind()) && isNumeric(inFnField.Kind()) {
		//	v := reflect.Zero(inFnField).Interface()
		//	setNumeric(args[i], &v)
		//	ins = append(ins, reflect.ValueOf(v))
		//	continue
		//}
		return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("go code error: fn in args type is not match"))
	}
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		return "", nil
	}
	if fnt.Out(len(outs) - 1).Implements(reflect.ValueOf(new(error)).Elem().Type()) {
		if errValueOf := outs[len(outs)-1]; !errValueOf.IsNil() {
			log.ZError(ctx, "fn call error", errValueOf.Interface().(error), "function name", funcName, "cost time", time.Since(t))
			return nil, errValueOf.Interface().(error)
		}
		if len(outs) == 1 {
			return "", nil
		}
		outs = outs[:len(outs)-1]
	}
	for i := 0; i < len(outs); i++ {
		out := outs[i]
		switch out.Kind() {
		case reflect.Map:
			if out.IsNil() {
				outs[i] = reflect.MakeMap(out.Type())
			}
		case reflect.Slice:
			if out.IsNil() {
				outs[i] = reflect.MakeSlice(out.Type(), 0, 0)
			}
		}
	}
	if len(outs) == 1 {
		log.ZInfo(ctx, "output resp", "function name", funcName, "resp", outs[0].Interface(), "cost time", time.Since(t))
		return outs[0].Interface(), nil
	}
	val := make([]any, 0, len(outs))
	for i := range outs {
		val = append(val, outs[i].Interface())
	}
	log.ZInfo(ctx, "output resp", "function name", funcName, "resp", val, "cost time", time.Since(t))
	return val, nil
}

// call 异步地调用指定的函数，并通过回调接口返回执行结果或错误信息。
// 参数：
//   - callback: 回调接口，提供成功和错误的回调方法。
//   - operationID: 调用操作的唯一标识符，用于追踪和日志记录。
//   - fn: 需要被调用的函数指针。
//   - args: 可变参数列表，传递给被调用函数的参数。

// 函数执行流程：
// 1. 检查 `callback` 是否为 `nil`，如果是，记录警告日志并返回。
// 2. 启动一个新的 Goroutine 来执行异步调用：
//   - 调用 `call_` 函数执行 `fn` 并获取结果 `res` 和错误 `err`。
//   - 如果发生错误，检查错误类型是否实现了 `CodeError` 接口：
//   - 如果是 `CodeError` 类型，调用 `callback.OnError` 并传递错误码和错误信息。
//   - 否则，将错误标记为 `UnknownCode` 并调用 `callback.OnError`。
//   - 如果没有错误，将 `res` 序列化为 JSON 字符串；若序列化失败，调用 `callback.OnError` 并返回内部错误信息。
//   - 如果成功序列化，将结果通过 `callback.OnSuccess` 返回给调用方。
func call(callback sdk_callback.Base, operationID string, fn any, args ...any) {
	if callback == nil {
		log.ZWarn(context.Background(), "callback is nil", nil)
		return
	}
	go func() {
		res, err := call_(operationID, fn, args...)
		if err != nil {
			if code, ok := err.(sdkerrs.CodeError); ok {
				callback.OnError(int32(code.Code()), code.Error())
			} else {
				callback.OnError(sdkerrs.UnknownCode, fmt.Sprintf("error %T not implement CodeError: %s", err, err))
			}
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			callback.OnError(sdkerrs.SdkInternalError, fmt.Sprintf("function res json.Marshal error: %s", err))
			return
		}
		callback.OnSuccess(string(data))
	}()
}

// syncCall 同步调用一个函数，并将其结果转换为 JSON 字符串返回。
// 参数：
//   - operationID: 调用操作的唯一标识符。
//   - fn: 需要被调用的函数指针。
//   - args: 可变参数列表，传递给被调用函数的参数。
// 返回值：
//   - string: 被调用函数的结果，序列化为 JSON 格式的字符串。

// 函数执行流程：
// 1. 检查 `operationID` 是否为空，若为空则返回空字符串。
// 2. 检查资源是否已加载，未加载则返回空字符串。
// 3. 通过反射获取 `fn` 的类型信息，验证 `fn` 是否为函数类型。
// 4. 检查 `fn` 的参数数量是否与传入的 `args` 数量匹配。
// 5. 准备实际调用的参数列表 `ins`：
//   - 创建上下文 `ctx` 并将其作为第一个参数。
//   - 对 `args` 进行类型检查和转换（如 JSON 字符串解析或数值类型转换），确保其与 `fn` 的参数类型兼容。
//
// 6. 使用反射调用 `fn`，获取其返回值 `outs`。
// 7. 检查函数的最后一个返回值是否为 `error` 类型：
//   - 如果存在且不为 `nil`，则记录错误日志并返回空字符串。
//   - 如果仅有一个返回值且为 `error`，则返回空字符串。
//
// 8. 将返回值中的 `nil` 类型（如 `map` 和 `slice`）转换为相应的空值（非 `nil`）。
// 9. 将返回值序列化为 JSON 格式的字符串；如果序列化失败，则返回空字符串。
// 10. 记录函数调用的输出日志，并返回序列化后的 JSON 字符串。
func syncCall(operationID string, fn any, args ...any) string {
	if operationID == "" {
		//callback.OnError(constant.ErrArgs.ErrCode, errs.ErrArgs.WrapMsg("operationID is empty").Error())
		return ""
	}
	if err := CheckResourceLoad(UserForSDK, ""); err != nil {
		return ""
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		//callback.OnError(10000, "go code error: fn is not function")
		return ""
	}
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	if len(args)+1 != numIn {
		//callback.OnError(10000, "go code error: fn in args num is not match")
		return ""
	}
	ins := make([]reflect.Value, 0, numIn)

	ctx := ccontext.WithOperationID(UserForSDK.BaseCtx(), operationID)
	t := time.Now()
	log.ZInfo(ctx, "input req", "function name", funcName, "args", args)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		tag := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
				v := reflect.New(tag)
				if args[i].(string) != "" {
					if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
						log.ZWarn(ctx, "json.Unmarshal error", err, "function name", funcName, "arg", args[i], "v", v.Interface())
						//callback.OnError(constant.ErrArgs.ErrCode, err.Error())
						return ""
					}
				}

				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		//callback.OnError(constant.ErrArgs.ErrCode, "go code error: fn in args type is not match")
		return ""
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.TypeOf((*error)(nil)).Elem())
	}
	fmt.Println("fnv:", fnv.Interface(), "ins:", ins)
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		//callback.OnSuccess("")
		return ""
	}
	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			//callback.OnError(10000, last.(error).Error())
			return ""
		}
		if len(outs) == 1 {
			//callback.OnSuccess("") // 只有一个返回值为error，且error == nil
			return ""
		}
		outVals = outVals[:len(outVals)-1]
	}
	// 将map和slice的nil转换为非nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		//callback.OnError(constant.ErrArgs.ErrCode, err.Error())
		return ""
	}
	log.ZInfo(ctx, "output resp", "function name", funcName, "resp", jsonVal, "cost time", time.Since(t))
	return string(jsonData)
}

// messageCall 异步地调用指定的函数 `fn`，并通过回调接口 `callback` 返回执行结果或错误信息。
// 参数：
//   - callback: 回调接口，提供成功和错误的回调方法。
//   - operationID: 操作的唯一标识符，用于追踪和日志记录。
//   - fn: 需要被调用的函数。
//   - args: 可变参数列表，传递给被调用函数的参数。

// 函数执行流程：
// 1. 检查 `callback` 是否为 `nil`，如果是，记录警告日志并返回。
// 2. 启动一个新的 Goroutine 来执行 `messageCall_`，以便异步处理函数调用和结果返回。

// messageCall_ 实际执行函数调用，并通过回调接口处理函数的执行结果或错误信息。
// 参数和功能与 `messageCall` 类似，只是它在一个新 Goroutine 中运行。
//
// 具体执行流程：
// 1. 使用 `defer` 和 `recover` 捕获运行时 panic，避免程序崩溃，并通过回调返回错误信息。
// 2. 验证 `operationID` 是否为空，如果为空，返回错误。
// 3. 检查资源是否已加载，若未加载，返回资源加载错误。
// 4. 确保 `fn` 为函数类型，否则返回内部错误。
// 5. 检查传入的参数数量是否与函数签名匹配，不匹配则返回错误。
// 6. 创建调用参数列表 `ins`，将 `ctx` 和函数参数类型转换后依次添加到列表中。
// 7. 根据反射调用函数 `fn`，获取返回值 `outs`。
// 8. 检查返回值中是否包含错误信息，如果有，则调用 `OnError` 回调并返回。
// 9. 将 `nil` 的 `map` 和 `slice` 转换为非 `nil` 的空集合。
// 10. 将返回值序列化为 JSON 格式，如果序列化成功，调用 `OnSuccess` 回调返回结果。

func messageCall(callback sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	if callback == nil {
		log.ZWarn(context.Background(), "callback is nil", nil)
		return
	}
	go messageCall_(callback, operationID, fn, args...)
}
func messageCall_(callback sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" panic err:", r, string(debug.Stack()))
			callback.OnError(sdkerrs.SdkInternalError, fmt.Sprintf("recover: %+v", r))
			return
		}
	}()
	if operationID == "" {
		callback.OnError(sdkerrs.ArgsError, sdkerrs.ErrArgs.WrapMsg("operationID is empty").Error())
		return
	}
	if err := CheckResourceLoad(UserForSDK, ""); err != nil {
		callback.OnError(sdkerrs.ResourceLoadNotCompleteError, "resource load error: "+err.Error())
		return
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		callback.OnError(sdkerrs.SdkInternalError, "go code error: fn is not function")
		return
	}
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	fmt.Println("fn args num is", numIn, len(args))
	if len(args)+1 != numIn {
		callback.OnError(sdkerrs.SdkInternalError, "go code error: fn in args num is not match")
		return
	}
	ins := make([]reflect.Value, 0, numIn)
	ctx := ccontext.WithOperationID(UserForSDK.BaseCtx(), operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, callback)

	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ { // callback sdk_callback.Base, operationID string, ...
		tag := fnt.In(i + 1) // ctx context.Context, ...
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
				v := reflect.New(tag)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					callback.OnError(sdkerrs.ArgsError, err.Error())
					return
				}
				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		callback.OnError(sdkerrs.ArgsError, "go code error: fn in args type is not match")
		return
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.ValueOf(new(error)).Elem().Type())
	}
	//fmt.Println("fnv:", fnv.Interface(), "ins:", ins)
	outs := fnv.Call(ins)

	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			if code, ok := last.(error).(errs.CodeError); ok {
				callback.OnError(int32(code.Code()), code.Error())
			} else {
				callback.OnError(sdkerrs.UnknownCode, fmt.Sprintf("error %T not implement CodeError: %s", last.(error), last.(error).Error()))
			}
			return
		}

		outVals = outVals[:len(outVals)-1]
	}
	// 将map和slice的nil转换为非nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		callback.OnError(sdkerrs.ArgsError, err.Error())
		return
	}
	callback.OnSuccess(string(jsonData))
}
