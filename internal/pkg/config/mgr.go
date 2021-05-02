package config

import "fmt"

func init() {
	fmt.Println("Initialization...")
	appConfig()
	ferretJobConfigs()
	newScheduler()
}
