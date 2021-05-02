/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"crawl/internal/pkg/config"
	"crawl/internal/pkg/fql"
	"crawl/internal/pkg/schedule"
)

func main() {
	// cmd.Execute()

	getJobs()
	<-make(chan bool)
}

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func getJobs() {

	ferret := fql.GetFerret()
	jobs := config.GetFerretJobs()

	// compile all fql scripts
	// the compiled program saved in Ferret struct
	// use ferret.ExecuteProgram(query key string) to run ferret query
	ferret.CompileAll(jobs)

	// execute all ferret queries
	//
	// 经测试，在for循环中，当调用Do时，传入的jobFunc函数名相同时，会被最后一个覆盖，
	// 假设 for 循环10次，那么相当于最后一次所设置的Do函数，会被重复执行10次，这时可加上一个『值不同』的参数来差异化
	// 这个值可以是数值，也可以是指针地址的值
	for _, j := range *jobs {
		// because gocron.Scheduler.Do in for-loop,
		// jobFunc (also with Pointer-type params is the same) will be overwrited, until the last trip
		//
		// In this loop, j is temporarily data copy, any change would lost,
		// unless j is a pointer itself which scope is outside the for-loop.
		//
		// schedule.MakeSchedule().Do(ferret.ExecuteProgramAndSaveOutput, args)，
		// make sure *args* is a value but not a temporary pointer,
		// even value of pointer will work too.
		//
		schedule.MakeSchedule(j).Do(ferret.ExecuteProgramAndSaveOutput, j)
	}
	schedule.Start(int(config.GetConfig().Scheduler.Max))
}
