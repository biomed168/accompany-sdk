package sdkerrs

var (
	ErrArgs           = NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")

	ErrResourceLoad = NewCodeError(ResourceLoadNotCompleteError, "ResourceLoadNotCompleteError")
)
