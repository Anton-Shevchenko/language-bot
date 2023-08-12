package jobManager

import (
	"time"
)

func (j *JobManager) Scheduler() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				j.do(j.findCurrentTasks(j.pull))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (j *JobManager) do(jobs []*Job) {
	for _, job := range jobs {
		go job.jobFunc()
	}
}

func (j *JobManager) findCurrentTasks(pull jobPull) []*Job {
	var tasks []*Job
	intervals := j.getPullIntervals(pull)

	for _, interval := range intervals {
		if j.isCurrentTask(interval) {
			tasks = append(tasks, pull[interval]...)
		}
	}

	return tasks
}

func (j *JobManager) isCurrentTask(interval int) bool {
	return interval%time.Now().Minute() != 0
}

func (j *JobManager) getPullIntervals(pull jobPull) []int {
	keys := make([]int, len(pull))

	i := 0
	for k := range pull {
		keys[i] = k
		i++
	}

	return keys
}
