package schedule

import (
	"crawl/internal/pkg/fql"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler

// cant set SingletonMode, SetMaxConcurrentJobs here
// otherwise,  scheduler will hold
func newScheduler() *gocron.Scheduler {

	scheduler = gocron.NewScheduler(time.UTC)
	return scheduler
	// sch.SingletonMode()
	// SetMaxConcurrentJobs
	// sch.SetMaxConcurrentJobs(int(GetConfig().Scheduler.Max), gocron.WaitMode)

}

// no need initialization for caller
func GetScheduler() *gocron.Scheduler {
	if scheduler != nil {
		return scheduler
	} else {
		return newScheduler()
	}
}

// start the cron loop and go back to main thread
func Start(concurrent int) {
	fmt.Println("Start Cron Service with max concurrent ", concurrent)
	GetScheduler().SetMaxConcurrentJobs(concurrent, gocron.WaitMode)
	GetScheduler().StartAsync()
}

// start the cron loop and blocking the main thread
func StartBlocking(concurrent int) {
	Start(concurrent)
	<-make(chan bool)
}

// schedule fql job
func MakeSchedule(job *fql.Job) *gocron.Scheduler {

	var s *gocron.Scheduler
	switch {
	case len(job.Schedule.Cron) > 0:
		s = GetScheduler().Cron(job.Schedule.Cron)
	case len(job.Schedule.Every) > 0:
		s = GetScheduler().Every(job.Schedule.Every)
		if len(job.Schedule.At) > 0 {
			s.At(job.Schedule.At)
		}
	default:
		s = GetScheduler().Every("7m")
	}
	log.Printf("fql job '%s(%s)' at schedule: %s", job.Name, job.Desc, job.Schedule)
	return s.Tag(job.Key)
}
