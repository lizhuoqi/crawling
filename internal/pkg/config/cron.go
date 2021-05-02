package config

import (
	"crawl/internal/pkg/fql"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler

// cant set SingletonMode, SetMaxConcurrentJobs here
// otherwise,  scheduler will hold
func newScheduler() {

	scheduler = gocron.NewScheduler(time.UTC)
	// sch.SingletonMode()
	// SetMaxConcurrentJobs
	// sch.SetMaxConcurrentJobs(int(GetConfig().Scheduler.Max), gocron.WaitMode)

}

func GetScheduler() *gocron.Scheduler {
	return scheduler
}

// start the cron loop and go back to main thread
func RunCron() {
	fmt.Println("Start Cron Service.")
	scheduler.SetMaxConcurrentJobs(int(GetConfig().Scheduler.Max), gocron.WaitMode)
	scheduler.StartAsync()
}

// start the cron loop and blocking the main thread
func RunCronBlocking() {
	fmt.Println("Start Cron Service.")
	scheduler.SetMaxConcurrentJobs(int(GetConfig().Scheduler.Max), gocron.WaitMode)
	scheduler.StartBlocking()
}

// schedule fql job
func MakeSchedule(job *fql.Job) *gocron.Scheduler {

	var s *gocron.Scheduler
	switch {
	case len(job.Schedule.Cron) > 0:
		s = scheduler.Cron(job.Schedule.Cron)
	case len(job.Schedule.Every) > 0:
		s = scheduler.Every(job.Schedule.Every)
		if len(job.Schedule.At) > 0 {
			s.At(job.Schedule.At)
		}
	default:
		s = scheduler.Every("7m")
	}

	return s.Tag(job.Key)
}
