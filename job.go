package samovar

import (
	"sync"
	"time"
	//"log"
	"fmt"
	"math/rand"
)

type Job struct {
	Data     funcarg
	Title    string
	jonstart time.Time
	Id       string
	//Number of times which job was call
	numberofcalls int
	done          chan bool
	lock          *sync.Mutex
	options       Options
	//delay in seconds
	delay time.Duration
}

//This struct provides basic options for job
type Options struct {
	immutable bool
}

//CreateJson provides creation of the new job
func CreateJob(title string, fn funcarg) *Job {
	job := &Job{}
	job.Data = fn
	job.Title = title
	job.Id = _generateid()
	job.done = make(chan bool)
	job.numberofcalls = 0
	job.delay = 0
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

//This method provides basic run of the job
func (j *Job) Run(arguments interface{}) {
	fmt.Println(fmt.Sprintf("Start to processing %s", j.Title))
	j.jobRun(arguments)
}

func (j *Job) IsDone() bool {
	return <-j.done
}

//Run current job with arguments
func (j *Job) jobRun(arguments interface{}) {
	go func() {
		time.Sleep(time.Second)
		j.Data(arguments)
		j.lock.Lock()
		j.numberofcalls++
		j.lock.Unlock()
		j.done <- true
	}()
}
