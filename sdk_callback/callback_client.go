package sdk_callback

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

type OnConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(errCode int32, errMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
}
