package fql

import "fmt"

type Job struct {
	Key    string
	Enable bool
	Name   string
	Desc   string
	Script string
	Output string
	// Compiled bool
	// binary
	Schedule struct {
		Every string
		At    string
		Cron  string
	}
}

func (job *Job) Run() (result []byte, err error) {
	return GetFerret().ExecuteProgram(job)
}

func (job *Job) RunAndSave() error {
	return GetFerret().ExecuteProgramAndSaveOutput(*job)
}

func (job *Job) Compile() error {
	return GetFerret().Compile(job)
}

// you have to initialize the variable of type Jobs
// e.g. `internal/pkg/config/job.go#init`
type Jobs map[string]Job

func (jobs Jobs) AddJob(job Job) {
	jobs[job.Key] = job
}

func (self Jobs) AddJobs(jobs []Job) {
	for _, j := range jobs {
		fmt.Println(j)
		self[j.Key] = j
	}
}

func (jobs Jobs) GetJob(name string) Job {
	job := jobs[name]
	return job
}

func (jobs Jobs) GetJobs() []Job {
	_jobs := make([]Job, 0)
	for _, v := range jobs {
		_jobs = append(_jobs, v)
	}

	return _jobs
}

func (jobs Jobs) Len() int {
	return len(jobs)
}
