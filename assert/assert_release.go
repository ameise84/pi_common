//go:build !debug

package assert

func Assert(f func() bool, printParam ...any) {
	funcName, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("===============================\nUnDel Assert\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
		file,
		line,
		runtime.FuncForPC(funcName).Name(),
		printParam,
	)
	panic(s)
}

func AssertNoError(f func() error, printParam ...any) {
	funcName, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("===============================\nUnDel Assert\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
		file,
		line,
		runtime.FuncForPC(funcName).Name(),
		printParam,
	)
	panic(s)
}

func EnSure(b bool, printParam ...any) {
	funcName, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("===============================\nUnDel EnSure\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
		file,
		line,
		runtime.FuncForPC(funcName).Name(),
		printParam,
	)
	panic(s)
}

func EnSureNoErr(e error, printParam ...any) {
	funcName, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("===============================\nUnDel EnSure\n>> File:\n\t%s:%d\n>> Func:%s\n>> Param:%+v\n===============================\n",
		file,
		line,
		runtime.FuncForPC(funcName).Name(),
		printParam,
	)
	panic(s)
}
