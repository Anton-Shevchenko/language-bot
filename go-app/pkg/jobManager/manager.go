package jobManager

import (
	"time"
)

type JobManager struct {
	pull jobPull
}

func NewJobManager() *JobManager {
	return &JobManager{
		pull: make(jobPull),
	}
}

type Job struct {
	jobFunc  jobFunc
	key      string
	interval time.Duration
}

type jobFunc func()
type jobPull map[int][]*Job

func (j *JobManager) Add(job jobFunc, key string, interval time.Duration) {
	minute := int(interval.Minutes())
	j.pull[minute] = append(j.pull[minute], &Job{
		key:     key,
		jobFunc: job,
	})
}
