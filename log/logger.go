package log

import (
	"github.com/ameise84/logger"
)

var _gLogger Logger

func init() {
	_gLogger = logger.DefaultLogger()
}
