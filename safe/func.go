package safe

import "runtime"

type PanicHook interface {
	OnPanic(err error)
}

// Func 安全的执行f函数,如果panic会自动捕获
func Func(h PanicHook, where string, f func()) {
	doFunc(h, where, f)
}

// LoopFunc 循环执行函数,直到其正常退出.通常用于func函数需要循环执行,被异常中止后需要继续执行
func LoopFunc(h PanicHook, where string, f func()) {
	for {
		if doFunc(h, where, f) {
			break
		}
		runtime.Gosched()
	}
}

func doFunc(h PanicHook, where string, f func()) bool {
	defer RecoverPanic(h, where)
	f()
	return true
}
