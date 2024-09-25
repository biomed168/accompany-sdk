package sdk_struct

import "accompany-sdk/ai_struct"

type SDKConfig struct {
	PlatformID           int32              `json:"platformID"`
	DataDir              string             `json:"dataDir"`
	LogLevel             uint32             `json:"logLevel"`
	IsLogStandardOutput  bool               `json:"isLogStandardOutput"`
	LogFilePath          string             `json:"logFilePath"`
	IsExternalExtensions bool               `json:"isExternalExtensions"`
	AiConfig             ai_struct.AiConfig `json:"aiConfig"`
}
