//go:build !debug

package safe

import (
	"errors"
	"fmt"
	"github.com/ameise84/pi_common/log"
	"runtime/debug"
)

func RecoverPanic(h PanicHook, where string) {
	if x := recover(); x != nil {
		stack := string(debug.Stack())
		msg := fmt.Sprintf("%+v[where:%s]\n stack:%s", x, where, stack)
		if h != nil {
			h.OnPanic(errors.New(msg))
		} else {
			log.Error(msg)
		}
	}
}
