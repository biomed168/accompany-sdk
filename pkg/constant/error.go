package constant

type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

func (e *ErrInfo) Error() string {
	return e.ErrMsg
}
