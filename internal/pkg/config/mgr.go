package config

import (
	"crawl/internal/pkg/schedule"
	"fmt"
)

func init() {
	fmt.Println("Initialization...")
	appConfig()
	ferretJobConfigs()
	schedule.NewScheduler()
}
