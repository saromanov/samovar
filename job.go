package samovar

import (
	"sync"
	"time"
	//"log"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

type Job struct {
	Data      interface{}
	Title     string
	Arguments interface{}
	jonstart  time.Time
	Id        string
	//Number of times which job was call
	numberofcalls int
	done          bool
	started       bool
	lock          *sync.RWMutex
	options       Options
	//delay in seconds
	delay         time.Duration
	result        interface{}
	executionTime float64
	executionTimes []float64
}

//This struct provides basic options for job
type Options struct {
	immutable bool
}

//CreateJson provides creation of the new job
func CreateJob(title string, fn interface{}) *Job {
	job := &Job{}
	job.Data = fn
	job.Title = title
	job.Id = _generateid()
	job.Arguments = []interface{}{}
	job.numberofcalls = 0
	job.delay = 0
	job.result = make(chan interface{})
	job.executionTimes = []float64{}
	job.lock = &sync.RWMutex{}
	return job
}

//Generate id for each job
func _generateid() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//RunWithDelay provides running job with delay n seconds
func (j *Job) RunWithDelay(arguments interface{}, delay uint) {
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		j.Run(arguments)
	}()
}

//RunEvery provides run job every n seconds
func (j *Job) RunEvery(arguments interface{}, sec uint) {
	go func() {
		for {
			time.Sleep(time.Duration(sec) * time.Second)
			j.Run(arguments)
		}
	}()
}

//RunTimes provides running of job t n times with delay sec
func (j *Job) RunTimes(arguments interface{}, sec uint, times int) {
	go func() {
		for i := 0; i < times; i++ {
			time.Sleep(time.Duration(sec) * time.Second)
			j.Run(arguments)
		}
	}()
}

//This method provides basic run of the job
func (j *Job) Run(arguments interface{}) {
	fmt.Println(fmt.Sprintf("Start to processing %s", j.Title))
	newarguments := j.prepareArguments(arguments)
	j.jobRun(newarguments)
}

//This method doing preparation arguments, before putting to the function
func (j *Job) prepareArguments(arguments interface{}) []reflect.Value {
	return []reflect.Value{reflect.ValueOf(arguments)}
}

//IsDone provides checking of current job is done
func (j *Job) IsDone() bool {
	if j.started {
		return false
	}
	return j.done
}

//AvgExecutionTime returns average time of execution of past results
func (j *Job) AvgExecutionTime()float64 {
	pastres := len(j.executionTimes)
	if pastres == 0 {
		return 0
	}
	result := 0.0
	for _, item := range j.executionTimes {
		result += item
	}

	return result/float64(pastres)
}

//Run current job with arguments
func (j *Job) jobRun(arguments []reflect.Value) {
	j.lock.RLock()
	j.numberofcalls++
	j.lock.RUnlock()
	go func() {
		starttime := time.Now()
		//j.started <- true
		j.done = false
		result := reflect.ValueOf(j.Data).Call(arguments)
		if len(result) > 0 {
			j.result = result[0].Interface()
		} else {
			j.result = nil
		}
		//j.started <- false
		j.done = true
		j.executionTime = time.Now().Sub(starttime).Seconds()
		j.executionTimes = append(j.executionTimes, j.executionTime)
		j.lock.Lock()
		j.lock.Unlock()
	}()
}

//getResult provides getting result from current job
func (j *Job) getResult() (interface{}, error) {
	if j.result == nil {
		return nil, errors.New("Task not contain return value")
	}
	return j.result, nil
}

func (j *Job) waitUntilResult() (interface{}, error) {
	for {
		if j.done == true {
			if j.result == nil {
				return nil, errors.New("Task not contain return value")
			}
			return j.result, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}
