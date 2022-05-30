package signal

import (
	"sync"

	"github.com/chuan-fu/Common/zlog"
)

type shutDown struct {
	fns []func() error
	sync.Mutex
}

func NewShutDown() *shutDown {
	return &shutDown{
		fns: make([]func() error, 0),
	}
}

func (s *shutDown) Add(f func() error) {
	s.Lock()
	defer s.Unlock()
	s.fns = append(s.fns, f)
}

func (s *shutDown) Run() {
	s.Lock()
	defer s.Unlock()
	for _, f := range s.fns {
		if err := f(); err != nil {
			log.Error("shutdown func error => ", err.Error())
		}
	}
}
