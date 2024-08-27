package errors

import (
	"errors"
	"fmt"
	"runtime"
)


func NewOrWrapNoStack(err error, msg string) error {
	if err == nil {
		return New(msg)
	} else {
		return Wrap(err, msg)
	}
}

func NewNoStack(msg string) error {
	return errors.New(msg)
}

func NewNoStackPrintf(format string, a ...any) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func WrapNoStack(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.New(msg + " -> " + err.Error())
}

func WrapNoStackPrintf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	return errors.New(fmt.Sprintf(format, a...) + " -> " + err.Error())
}

func NewOrWrap(err error, msg string) error {
	//return newOrWrapWithStack(err, msg,2)
	return NewOrWrapNoStack(err, msg)
}

func NewOrWrapPrintf(err error, format string, a ...any) error {
	//return newOrWrapWithStack(err, fmt.Sprintf(format, a...),2)
	return NewOrWrapNoStack(err,fmt.Sprintf(format, a...))
}

func New(msg string) error {
	//return newWithStack(msg, 2)
	return  NewNoStack(msg)
}

func NewPrintf(format string, a ...any) error {
	//return newWithStack(fmt.Sprintf(format, a...),2)
	return NewNoStack(fmt.Sprintf(format, a...))
}

func Wrap(err error, msg string) error {
	//return wrapWithStack(err, msg,2)
	return WrapNoStack(err,msg)
}

func WrapPrintf(err error, format string, a ...any) error {
	//return wrapWithStack(err, fmt.Sprintf(format, a...),2)
	return WrapNoStack(err,fmt.Sprintf(format, a...))
}

func NewOrWrapWithStack(err error, msg string) error {
	return newOrWrapWithStack(err, msg, 2)
}

func NewWithStack(msg string) error {
	return newWithStack(msg, 2)
}

func NewWithStackPrintf(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return newWithStack(msg,2)
}

func WrapWithStack(err error, msg string) error {
	return wrapWithStack(err,msg,2)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func newOrWrapWithStack(err error, msg string, skip int) error {
	if err == nil {
		return newWithStack(msg,skip)
	} else {
		return wrapWithStack(err, msg,skip)
	}
}

func newWithStack(msg string, skip int) error {
	funcName, file, line, _ := runtime.Caller(skip)
	format := "func: %s\t%s:%d"
	stack := fmt.Sprintf(format, runtime.FuncForPC(funcName).Name(), file, line)
	return &stackError{msg: msg, stack: stack}
}

func wrapWithStack(err error, msg string,  skip int) error {
	if err == nil {
		return nil
	}
	funcName, file, line, _ := runtime.Caller(skip)
	format := "func: %s\t%s:%d"
	stack := fmt.Sprintf(format, runtime.FuncForPC(funcName).Name(), file, line)
	var e *stackError
	switch {
	case errors.As(err, &e):
		return &stackError{
			msg:   msg + " -> " + e.msg,
			stack: stack + "\n" + e.stack,
		}
	default:
		return &stackError{
			msg:   msg + " -> " + err.Error(),
			stack: stack,
		}
	}
}

