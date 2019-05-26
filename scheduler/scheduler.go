package scheduler

import (
	"sync/atomic"
)

// ShouldParallelize returns true iff there is a rasoneable number of slaves working
func ShouldParallelize() bool {
	return trafficWarden.shouldParallelize()
}

// AddJob must be called before starting to run a new slave
func AddJob() {
	trafficWarden.addJob()
}

// AddJobs must be called before starting to run numberofJobs a new slaves
func AddJobs(numberOfJobs int32) {
	trafficWarden.addJobs(numberOfJobs)
}

// JobEnded must be called after a slave has completed his job
func JobEnded() {
	trafficWarden.jobEnded()
}

// var trafficWarden = newScheduler(int32(runtime.NumCPU()))

var trafficWarden = newScheduler(6400)

type scheduler struct {
	slaves    int32
	maxSlaves int32
}

func newScheduler(maxSlaves int32) *scheduler {
	return &scheduler{
		slaves:    0,
		maxSlaves: maxSlaves,
	}
}

func (s *scheduler) shouldParallelize() bool {
	return atomic.LoadInt32(&(s.slaves)) <= s.maxSlaves
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
