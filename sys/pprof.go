package sys

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func DebugHttpPporf(port uint16, exit <-chan struct{}) {
	go func() {
		srv := http.Server{
			Addr: fmt.Sprintf(":%d", port),
		}
		go func() {
			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				_gLogger.ErrorPrintf("pprof http start: <%v>", err)
			}
		}()
		_ = <-exit
		if err := srv.Shutdown(context.Background()); err != nil {
			_gLogger.ErrorPrintf("pprof http stop: <%v>", err)
		}
	}()
}

func DebugFilePporf(path string, exit <-chan struct{}) {
	go func() {
		err := os.MkdirAll(path, 0666)
		if err != nil {
			_gLogger.Error(err)
			return
		}
		fc, err := os.Create(filepath.Join(path, "cpu.prof"))
		if err != nil {
			_gLogger.ErrorPrintf("pprof create cpu.prof: <%v>", err)
			return
		}
		fm, err := os.Create(filepath.Join(path, "memory.prof"))
		if err != nil {
			_gLogger.ErrorPrintf("pprof create memory.prof: <%v>", err)
			return
		}
		err = pprof.StartCPUProfile(fc)
		if err != nil {
			_gLogger.ErrorPrintf("pprof start cpu file: <%v>", err)
			return
		}
		_ = <-exit
		err = pprof.WriteHeapProfile(fm)
		if err != nil {
			_gLogger.ErrorPrintf("pprof write memory file: <%v>", err)
			return
		}
		pprof.StopCPUProfile()
	}()
}
