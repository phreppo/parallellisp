package cell

import (
	"runtime"
	"sync/atomic"
)

type scheduler struct {
	slaves    int32
	maxSlaves int32
}

func newScheduler() *scheduler {
	return &scheduler{
		slaves:    0,
		maxSlaves: int32(runtime.NumCPU()),
	}
}

func (s *scheduler) shouldParallelize() bool {
	return atomic.LoadInt32(&(s.slaves)) < s.maxSlaves
}

func (s *scheduler) addJob() {
	atomic.AddInt32(&(s.slaves), 1)
}

func (s *scheduler) addJobs(numberOfJobs int32) {
	atomic.AddInt32(&(s.slaves), numberOfJobs)
}

func (s *scheduler) jobEnded() {
	atomic.AddInt32(&(s.slaves), -1)
}
