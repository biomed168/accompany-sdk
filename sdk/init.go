package sdk

import (
	"accompany-sdk/sdk_callback"
	"accompany-sdk/sdk_struct"
	"encoding/json"
	"fmt"

	"github.com/openimsdk/tools/log"
)

const (
	rotateCount  uint = 1
	rotationTime uint = 24
	version           = "1.0"
)

func InitSDK(listener sdk_callback.OnConnListener, operationID string, config string) bool {
	var configArgs sdk_struct.SDKConfig
	if err := json.Unmarshal([]byte(config), &configArgs); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if err := log.InitFromConfig("sdk", "", int(configArgs.LogLevel), configArgs.IsLogStandardOutput, false, configArgs.LogFilePath, rotateCount, rotationTime, version, true); err != nil {
		fmt.Println(operationID, "log init failed ", err.Error())
	}
	fmt.Println("InitSDK success")
	UserForSDK = new(LoginMgr)
	return UserForSDK.InitSDK(configArgs, listener)
}

func Login(callback sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, UserForSDK.Login, userID, token)
}
