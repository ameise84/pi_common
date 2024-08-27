package common

import (
	"errors"
	"runtime"
	"sync/atomic"
)

const (
	stopped int32 = iota
	starting
	started
	stopping
)

type Service struct {
	state int32
}

func (s *Service) Start(f func() error) (err error) {
	if !atomic.CompareAndSwapInt32(&s.state, stopped, starting) {
		return errors.New("the cluster has not stopped")
	}
	if f != nil {
		if err = f(); err != nil {
			atomic.StoreInt32(&s.state, stopped)
		} else {
			atomic.StoreInt32(&s.state, started)
		}
	} else {
		atomic.StoreInt32(&s.state, started)
	}
	return
}

func (s *Service) Stop(f func()) {
	for {
		if atomic.LoadInt32(&s.state) == stopped {
			break
		}
		if atomic.CompareAndSwapInt32(&s.state, started, stopping) {
			if f != nil {
				f()
			}
			atomic.StoreInt32(&s.state, stopped)
			break
		}
		runtime.Gosched()
	}
}

func (s *Service) IsRunning() bool {
	return atomic.LoadInt32(&s.state) == started

}

func (s *Service) IsStopped() bool {
	return atomic.LoadInt32(&s.state) == stopped
}
