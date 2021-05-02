package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	JobsDir   string // default configDir/jobs
	FqlDir    string // default configDir/fql
	OutDir    string // default current working directory
	Scheduler struct {
		Max uint8 // maximum fql crawling jobs running cocurrenly running
	}

	Fql struct {
		Cache bool // cache fql compiled programm in memory
	}
}

var cfgFile string
var c Config

// appConfig reads in config file and ENV variables if set.
func appConfig() {
	setDefault()
	// todo
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// todo default settings
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigName("crawl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Sprintf("Unable to decode into Config struct , %v", err))
	}

	// make the `outdir` directory if not exists, and check
	// c.OutDir will never nil , beacuase setDefault()
	if err := os.Mkdir(c.OutDir, 0755); !(err == nil || os.IsExist(err)) {
		fmt.Printf("Can't make directory %s, because %v", c.OutDir, err)
		// fall back to current working directory "."
		c.OutDir = "."
	}
}

func GetConfig() *Config {
	return &c
}

func setDefault() {
	viper.SetDefault("JobsDir", "configDir/jobs")
	viper.SetDefault("FqlDir", "configDir/fql")
	viper.SetDefault("OutDir", "output")
	viper.SetDefault("Scheduler", map[string]string{"max": "8"})
}
