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

	// compile all fql scripts
	// the compiled program saved in Ferret struct
	// use ferret.ExecuteProgram(query key string) to run ferret query
	ferret.CompileAll(config.GetFerretJobs())

	// execute all ferret queriesß
	for _, j := range *config.GetFerretJobs() {
		//  todo 结果再处理（保存，或重定向到文件）
		// 经测试，在for循环中，当调用Do时，传入的函数名相同时，会被最后一个覆盖，
		// 假设 for 循环10次，那么相当于最后一次所设置的Do函数，会被重复执行10次，这时可加上不同的参数来差异化
		schedule.MakeSchedule(&j).Do(ferret.ExecuteProgramAndSaveOutput, j)
	}
	schedule.RunCron(int(config.GetConfig().Scheduler.Max))
}
