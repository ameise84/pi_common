package log

type Logger interface {
	Trace(any)
	Debug(any)
	Info(any)
	Warn(any)
	Error(any)
	Fatal(any)
}

func SetLogger(l Logger) {
	_gLogger = l
}

func Trace(a any) {
	_gLogger.Trace(a)
}

func Debug(a any) {
	_gLogger.Debug(a)
}

func Info(a any) {
	_gLogger.Info(a)
}

func Error(a any) {
	_gLogger.Error(a)
}

func Fatal(a any) {
	_gLogger.Fatal(a)
}
