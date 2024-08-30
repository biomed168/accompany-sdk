package sdk

import (
	"accompany-sdk/pkg/ccontext"
	"accompany-sdk/pkg/sdkerrs"
	"accompany-sdk/sdk_callback"
	"accompany-sdk/sdk_struct"
	"context"
	"os/user"
	"strings"
	"sync"
)

const (
	LogoutStatus = iota + 1
	Logging
	Logged
)

var (
	UserForSDK *LoginMgr
)

// CheckResourceLoad checks the SDK is resource load status.
func CheckResourceLoad(uSDK *LoginMgr, funcName string) error {
	if uSDK == nil {
		return sdkerrs.New("SDK not initialized,userForSDK is nil", "funcName", funcName).Wrap()
	}
	if funcName == "" {
		return nil
	}
	parts := strings.Split(funcName, ".")
	if parts[len(parts)-1] == "Login-fm" {
		return nil
	}
	if uSDK.getLoginStatus(context.Background()) != Logged {
		return sdkerrs.New("SDK not logged in", "funcName", funcName).Wrap()
	}
	return nil
}

type LoginMgr struct {
	user         *user.User
	justOnceFlag bool

	w           sync.Mutex
	loginStatus int

	connListener sdk_callback.OnConnListener

	ctx       context.Context
	cancel    context.CancelFunc
	info      *ccontext.GlobalConfig
	id2MinSeq map[string]int64
}

func (u *LoginMgr) getLoginStatus(_ context.Context) int {
	u.w.Lock()
	defer u.w.Unlock()
	return u.loginStatus
}

func (u *LoginMgr) setLoginStatus(status int) {
	u.w.Lock()
	defer u.w.Unlock()
	u.loginStatus = status
}

func (u *LoginMgr) BaseCtx() context.Context {
	return u.ctx
}

func (u *LoginMgr) InitSDK(config sdk_struct.SDKConfig, listener sdk_callback.OnConnListener) bool {
	if listener == nil {
		return false
	}
	u.info = &ccontext.GlobalConfig{}
	u.info.SDKConfig = config
	u.connListener = listener
	return true
}

func (u *LoginMgr) login(ctx context.Context, userID, token string) error {
	u.setLoginStatus(Logged)
	return nil
}
