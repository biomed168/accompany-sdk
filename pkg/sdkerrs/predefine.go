package sdkerrs

import "github.com/openimsdk/tools/errs"

var (
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = errs.NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")

	ErrResourceLoad = errs.NewCodeError(ResourceLoadNotCompleteError, "ResourceLoadNotCompleteError")
)
