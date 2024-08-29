package sdk_struct

type MessageReceipt struct {
	GroupID     string   `json:"groupID"`
	UserID      string   `json:"userID"`
	MsgIDList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
}
type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         int64  `json:"seq"`
	Ex                          string `json:"ex"`
}
type MessageReaction struct {
	ClientMsgID  string `json:"clientMsgID"`
	ReactionType int    `json:"reactionType"`
	Counter      int32  `json:"counter,omitempty"`
	UserID       string `json:"userID"`
	GroupID      string `json:"groupID"`
	SessionType  int32  `json:"sessionType"`
	Info         string `json:"info,omitempty"`
}
type ImageInfo struct {
	Width  int32  `json:"x"`
	Height int32  `json:"y"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
}
type PictureBaseInfo struct {
	UUID   string `json:"uuid,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Url    string `json:"url,omitempty"`
}
type SoundBaseInfo struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
	SoundType string `json:"soundType,omitempty"`
}
type VideoBaseInfo struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize"`
	Duration       int64  `json:"duration"`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize"`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth"`
	SnapshotHeight int32  `json:"snapshotHeight"`
	SnapshotType   string `json:"snapshotType,omitempty"`
}
type FileBaseInfo struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize"`
	FileType  string `json:"fileType,omitempty"`
}

type SDKConfig struct {
	PlatformID           int32  `json:"platformID"`
	DataDir              string `json:"dataDir"`
	LogLevel             uint32 `json:"logLevel"`
	IsLogStandardOutput  bool   `json:"isLogStandardOutput"`
	LogFilePath          string `json:"logFilePath"`
	IsExternalExtensions bool   `json:"isExternalExtensions"`
}
