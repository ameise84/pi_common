package errors

type stackError struct {
	msg   string
	stack string
}

func (e *stackError) Error() string {
	return e.msg + "\n" + e.stack
}
