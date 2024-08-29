package sdk

import (
	"accompany-sdk/sdk_callback"
	"accompany-sdk/sdk_struct"
	"encoding/json"
	"fmt"
)

func InitSDK(listener sdk_callback.OnConnListener, operationID string, config string) bool {
	var configArgs sdk_struct.SDKConfig
	if err := json.Unmarshal([]byte(config), &configArgs); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	UserForSDK = new(LoginMgr)
	return UserForSDK.InitSDK(configArgs, listener)
}

func Login(callback sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, UserForSDK.login, userID, token)
}
