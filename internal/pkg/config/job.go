package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"crawl/internal/pkg/fql"
)

var jobs fql.Jobs

func init() {
	// prevent assignment to entry in nil map
	jobs = make(fql.Jobs)
}

func GetFerretJobs() *fql.Jobs {
	return &jobs
}

func ferretJobConfigs() {

	// get jobsdir value in crawl.yaml
	// viper is already initialized in `internal/pkg/config/app.go`
	jobsDir := viper.GetString("jobsdir")
	fqldir := viper.GetString("fqldir")
	outdir := viper.GetString("outdir")

	// scan job description yaml
	fmt.Printf("Scan jobs in directory %s. ", jobsDir)
	files, _ := walkThrough(jobsDir)

	fmt.Printf("%d files fo Job-Description founded.", len(files))

	// use fql.Job for Unmarshalling
	type jobsFromYaml struct {
		Enable  bool
		Fqljobs []fql.Job
	}
	// Get jobs in every yaml/json
	for _, f := range files {
		jobviper := viper.New()
		key := strings.TrimLeft(f[:len(f)-len(filepath.Ext(f))], jobsDir)
		// read fql job yaml
		jobviper.SetConfigFile(f)
		jobviper.ReadInConfig()

		var fromYaml jobsFromYaml
		jobviper.Unmarshal(&fromYaml)

		// if all jobs all disable
		if !fromYaml.Enable {
			continue
		}
		// turn to fql.job
		// modify Job.Key, Job.Script
		for _, j := range fromYaml.Fqljobs {
			// only deal with jobs of enable
			if j.Enable {
				// make should Key is unqiue for every fql.Job
				j.Key = key + "/" + j.Name
				// relative/absolute realpath of the script file
				j.Script = filepath.Join(fqldir, j.Script)
				// relative/absolute path of output file
				if len(j.Output) == 0 {
					// use script file
					j.Output = strings.Replace(j.Script, "/", "_", -1)
					j.Output = strings.Replace(j.Output, "\\", "_", -1)
					j.Output = j.Output[:len(j.Output)-len(filepath.Ext(j.Output))] + ".json"
				}
				j.Output = filepath.Join(outdir, j.Output)

				// make the `outdir` directory if not exists, and check
				// mkdir
				dir := filepath.Dir(j.Output)
				name := filepath.Base(j.Output)
				if err := os.Mkdir(dir, 0755); !(err == nil || os.IsExist(err)) {
					log.Printf("Can't make directory %s, because %v", dir, err)
					// fall back to crawl output dir
					j.Output = filepath.Join(c.OutDir, name)
				}
				// save job
				jobs.AddJob(j)
			}
		}
	}

	fmt.Printf("Eventually, get %d jobs.\n", jobs.Len())
}

//////// helper ///////
func walkThrough(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) == ".json" || filepath.Ext(info.Name()) == ".yaml" {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}
