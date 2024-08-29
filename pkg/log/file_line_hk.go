/*
get the name and line number of the calling file hook
*/
package log

import (
	"accompany-sdk/pkg/utils"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type fileHook struct{}

func newFileHook() *fileHook {
	return &fileHook{}
}

func (f *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//	funcation (f *fileHook) Fire(entry *logrus.Entry) error {
//		var s string
//		_, b, c, _ := runtime.Caller(8)
//		i := strings.SplitAfter(b, "/")
//		if len(i) > 3 {
//			s = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(c)
//		}
//		entry.Data["FilePath"] = s
//		return nil
//	}
func (f *fileHook) Fire(entry *logrus.Entry) error {
	var s string
	_, b, c, _ := runtime.Caller(8)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		s = b[i+1:len(b)] + ":" + utils.IntToString(c)
	}
	entry.Data["FilePath"] = s
	return nil
}
