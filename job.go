package samovar

import (
	"sync"
	"time"
	//"log"
	"fmt"
	"math/rand"
	"reflect"
)

type Job struct {
	Data     interface{}
	Title    string
	jonstart time.Time
	Id       string
	//Number of times which job was call
	numberofcalls int
	done          bool
	started       chan bool
	lock          *sync.Mutex
	options       Options
	//delay in seconds
	delay         time.Duration
	result        chan interface{}
	executionTime int64
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
	job.numberofcalls = 0
	job.delay = 0
	job.started = make(chan bool)
	job.result = make(chan interface{})
	job.lock = &sync.Mutex{}
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
func (j *Job) RunWithDelay(arguments []interface{}, delay uint) {
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		j.Run(arguments)
	}()
}

//RunEvery provides run job every n seconds
func (j *Job) RunEvery(arguments []interface{}, sec uint) {
	go func() {
		for {
			time.Sleep(time.Duration(sec) * time.Second)
			j.Run(arguments)
		}
	}()
}

//RunTimes provides running of job t n times with delay sec
func (j *Job) RunTimes(arguments []interface{}, sec uint, times int) {
	go func() {
		for i := 0; i < times; i++ {
			time.Sleep(time.Duration(sec) * time.Second)
			j.Run(arguments)
		}
	}()
}

//This method provides basic run of the job
func (j *Job) Run(arguments []interface{}) {
	fmt.Println(fmt.Sprintf("Start to processing %s", j.Title))
	newarguments := j.prepareArguments(arguments)
	j.jobRun(newarguments)
}

//THis method doing preparation arguments, before putting to the function
func (j *Job) prepareArguments(arguments []interface{})[]reflect.Value {
	res := []reflect.Value{}
	if len(arguments) == 0 {
		return res
	}
	for _, i := range arguments {
		res = append(res, reflect.ValueOf(i))
	}
	return res
}

//IsDone provides checking of current job is done
func (j *Job) IsDone() bool {
	return j.done
}

//Run current job with arguments
func (j *Job) jobRun(arguments []reflect.Value) {
	j.numberofcalls++
	go func() {
		starttime := time.Now().UnixNano()
		//j.started <- true
		j.result <- reflect.ValueOf(j.Data).Call(arguments)
		fmt.Println(<-j.result)
		//j.started <- false
		j.done = true
		j.executionTime = time.Now().UnixNano() - starttime
		j.lock.Lock()
		j.lock.Unlock()
	}()
}

func (j *Job) getResult() interface{} {
	return j.result
}
