package sys

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	_gOnce       sync.Once
	_gCancelFunc context.CancelFunc
	_gContext    context.Context
	_gExitLock   sync.Mutex
)

func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

var layout = "2006-01-02 15:04:05.000"
var traceFormat = "\033[96m[%s TRAC] %v\033[0m\n"
var errorFormat = "\033[91;1m[%s ERRO] %v\033[0m\n"

// WaitKillSigint 等待关闭信号
func WaitKillSigint() {
	_gOnce.Do(func() {
		_gContext, _gCancelFunc = context.WithCancel(context.Background())
	})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case _ = <-sigChan:
		fmt.Printf(traceFormat, time.Now().Format(layout), "接收到kill信号,执行清理中...")
	case <-_gContext.Done():
		fmt.Printf(traceFormat, time.Now().Format(layout), "服务器上下文中断,执行清理中...")
	}
}

func Exit(n int) {
	_gExitLock.Lock()
	defer _gExitLock.Unlock()
	s := fmt.Sprintf("进程异常关闭,原因[%d]", n)
	fmt.Printf(traceFormat, time.Now().Format(layout), s)
	if _gCancelFunc != nil {
		_gCancelFunc()
	} else {
		os.Exit(n)
	}
}
