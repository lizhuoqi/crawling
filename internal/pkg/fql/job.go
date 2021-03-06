package fql

import (
	"time"
)

type Job struct {
	Key    string
	Enable bool
	Name   string
	Desc   string
	Script string
	Output string
	// Compiled bool
	// binary
	Schedule Schedule

	Stats struct {
		LastStart time.Time
		LastStop  time.Time
		Times     int // sucessful times
		LastState JobState
		Duration  time.Duration // total time spend
	}
}
type Schedule struct {
	Every string
	At    string
	Cron  string
}

func (s Schedule) String() string {
	if s.Cron != "" {
		return s.Cron
	}
	if s.Every != "" {
		return s.Every
	}
	if s.At != "" {
		return s.At
	}
	return ""
}

// Not used yet.
// func (job *Job) Run() (result []byte, err error) {
// 	return GetFerret().ExecuteProgram(job)
// }

// func (job *Job) RunAndSave() error {
// 	return GetFerret().ExecuteProgramAndSaveOutput(job)
// }

// func (job *Job) Compile() error {
// 	return GetFerret().Compile(job)
// }

type JobRunner func(job *Job) (result []byte, err error)

func (job *Job) runnerMeasure(fun JobRunner) (result []byte, err error) {
	job.Stats.LastStart = time.Now()
	out, err := fun(job)
	job.Stats.LastStop = time.Now()
	if err != nil {
		job.Stats.LastState = Failed
	} else {
		// staticstic: only successful
		job.Stats.Duration += job.Stats.LastStop.Sub(job.Stats.LastStart)
		job.Stats.Times += 1
	}
	return out, err
}

// you have to initialize the variable of type Jobs
// e.g. `internal/pkg/config/job.go#init`
type Jobs map[string]*Job

func (jobs Jobs) AddJob(job Job) {
	jobs[job.Key] = &job
}

func (jobs Jobs) GetJob(name string) *Job {
	return jobs[name]

}

func (jobs Jobs) Len() int {
	return len(jobs)
}

// Not used yet
//
// func (self Jobs) AddJobs(jobs []Job) {
// 	for _, j := range jobs {
// 		self[j.Key] = &j
// 	}
// }
//
// func (jobs Jobs) GetJobs() []Job {
// 	_jobs := make([]Job, 0)
// 	for _, v := range jobs {
// 		_jobs = append(_jobs, v)
// 	}
//
// 	return _jobs
// }

////// job state enum /////
type JobState int8

const (
	Running JobState = iota
	Stopped
	Pending
	Terminated
	Finished
	Failed
	Unknown
)

func (state JobState) String() string {
	switch state {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Finished:
		return "Finished"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}
