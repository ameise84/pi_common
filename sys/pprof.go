package sys

import (
	"context"
	"errors"
	"fmt"
	"github.com/ameise84/pi_common/log"
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
				log.Error(fmt.Sprintf("pprof http start: <%v>", err))
			}
		}()
		_ = <-exit
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(fmt.Sprintf("pprof http stop: <%v>", err))
		}
	}()
}

func DebugFilePporf(path string, exit <-chan struct{}) {
	go func() {
		err := os.MkdirAll(path, 0666)
		if err != nil {
			log.Error(fmt.Sprintf("make dir: <%v>", err))
			return
		}
		fc, err := os.Create(filepath.Join(path, "cpu.prof"))
		if err != nil {
			log.Error(fmt.Sprintf("pprof create cpu.prof: <%v>", err))
			return
		}
		fm, err := os.Create(filepath.Join(path, "memory.prof"))
		if err != nil {
			log.Error(fmt.Sprintf("pprof create memory.prof: <%v>", err))
			return
		}
		err = pprof.StartCPUProfile(fc)
		if err != nil {
			log.Error(fmt.Sprintf("pprof start cpu file: <%v>", err))
			return
		}
		_ = <-exit
		err = pprof.WriteHeapProfile(fm)
		if err != nil {
			log.Error(fmt.Sprintf("pprof write memory file: <%v>", err))
			return
		}
		pprof.StopCPUProfile()
	}()
}
