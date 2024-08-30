// Copyright © 2024 OpenIM open source community. All rights reserved.
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

package sdkerrs

// 通用错误码
const (
	NetworkError             = 10000
	NetworkTimeoutError      = 10001
	ArgsError                = 10002 //输入参数错误
	CtxDeadlineExceededError = 10003 //上下文超时

	ResourceLoadNotCompleteError = 10004 //资源初始化未完成
	UnknownCode                  = 10005 //没有解析到code
	SdkInternalError             = 10006 //SDK内部错误
)

const (
	ServerInternalError = 500  // Server internal error
	NoPermissionError   = 1002 // Insufficient permission
	DuplicateKeyError   = 1003
	RecordNotFoundError = 1004 // Record does not exist

	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507
)
