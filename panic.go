package util

import (
	"fmt"
	"runtime"

	"github.com/archer-plus/util/logx"
)

// Recover 恢复panic
func Recover() {
	err := recover()
	if err != nil {
		logx.Sugar.Errorf("--RECOVER PANIC--: %v", err)
		logx.Sugar.Error(PrintStack())
	}
}

// PrintStack 打印Panic堆栈信息
func PrintStack() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
