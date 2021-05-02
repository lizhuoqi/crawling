package fql

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
)

// ferret Program wrapper
type Binary struct {
	Program *runtime.Program
}

func (bin *Binary) Run(ctx context.Context) (result []byte, err error) {
	return bin.Program.Run(ctx)
}

// put compiler/compiled(query program)/Context togather
type Ferret struct {
	Compiler *compiler.Compiler
	programs map[string]Binary
	Context  context.Context
}

func newFerret() *Ferret {
	ferret := Ferret{
		Compiler: compiler.New(),
		programs: make(map[string]Binary),
		Context:  defaultContext(),
	}

	return &ferret
}

// compile ferret query and cache the result
func (ferret *Ferret) Compile(job *Job) error {
	query, err := ioutil.ReadFile(job.Script)

	if err != nil {
		fmt.Printf("Read ferret script %s failed, %v.", job.Script, err)
		return err
	}

	program, err := ferret.Compiler.Compile(string(query))

	if err != nil {
		// return err
		fmt.Printf("Compile ferret script %s failed, %v.\n", job.Script, err)
		return err
	}

	fmt.Printf("Fql script %s compiled\n", job.Script)
	// update/add compiled fql script
	ferret.programs[job.Key] = Binary{Program: program}

	return nil
}

func (ferret *Ferret) CompileAll(jobs *Jobs) {
	for _, j := range *jobs {
		ferret.Compile(j)
	}
}

func (ferret *Ferret) GetCompiled(key string) *Binary {
	bin := ferret.programs[key]
	return &bin
}

func (ferret *Ferret) GetAllompiled() *(map[string]Binary) {
	return &ferret.programs
}

func (ferret *Ferret) execute(job *Job) (result []byte, err error) {
	fmt.Printf("Execute ferret query %s. %s. \n", job.Key, time.Now())
	bin := ferret.programs[job.Key]
	out, err := bin.Run(ferret.Context)

	if err != nil {
		fmt.Printf("fql script %s execute failed, %v.\n", job.Key, err)
		return nil, err
	} else {
		return out, nil
	}
}

func (ferret *Ferret) ExecuteProgram(job *Job) (result []byte, err error) {
	return job.runnerMeasure(ferret.execute)
}

func (ferret *Ferret) ExecuteProgramAndSaveOutput(job *Job) error {
	out, err := ferret.ExecuteProgram(job)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(job.Output, out, 0644); err != nil {
		fmt.Printf("Can't save output of query %s.\n", job.Key)
		return err
	}
	return nil
}

////////// helper /////////////

func defaultContext() context.Context {

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	// ctx = drivers.WithContext(ctx, http.NewDriver())
	// ctx = drivers.WithContext(ctx, cdp.NewDriver())
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())
	return ctx
}
