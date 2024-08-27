//go:build debug

//使用 build标签  -i -tags=debug 区分debug 和 release版本

package assert

import (
	"fmt"
	"runtime"
)

// Assert 只能模拟断言,在debug模式下会被求值
func Assert(f func() bool, printParam ...any) {
	if !f() {
		funcName, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("===============================\nAssert Fail\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
			file,
			line,
			runtime.FuncForPC(funcName).Name(),
			printParam,
		)
		panic(s)
	}
}

func AssertNoError(f func() error, printParam ...any) {
	if f() != nil {
		funcName, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("===============================\nAssert Fail\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
			file,
			line,
			runtime.FuncForPC(funcName).Name(),
			printParam,
		)
		panic(s)
	}
}

func EnSure(b bool, printParam ...any) {
	if !b {
		funcName, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("===============================\nAssert Fail\n>> File:%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
			file,
			line,
			runtime.FuncForPC(funcName).Name(),
			printParam,
		)
		panic(s)
	}
}

func EnSureNoErr(e error, printParam ...any) {
	if e != nil {
		funcName, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("===============================\nAssert Fail\n>> File:%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
			file,
			line,
			runtime.FuncForPC(funcName).Name(),
			printParam,
		)
		panic(s)
	}
}
