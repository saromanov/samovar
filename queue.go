package samovar

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//Options for queue
type QueueOptions struct {
}

//This queue provides basic data structure
type Queue struct {
	jobs        []*Job
	groupjobs   [][]*Job
	runningjobs int32
	title       string
	limit       int
	options     *QueueOptions
	dbstore     *redis.Client
}

func CreateQueue(title string) *Queue {
	queue := new(Queue)
	queue.jobs = []*Job{}
	queue.groupjobs = [][]*Job{}
	queue.title = title
	queue.runningjobs = 0
	//By default, queue is unlimited
	queue.limit = -1
	queue.options = new(QueueOptions)
	queue.dbstore = initRedis("localhost:6379")
	return queue
}

//This method provides getting and running job
func (q *Queue) Get() {

}

func (q *Queue) Put(job *Job, jp JobParams) {
	q.runJob(job, jp)
	q.jobs = append(q.jobs, job)
}

func (q *Queue) PutGroup(gjob []*Job) {
	q.groupjobs = append(q.groupjobs, gjob)
}

func (q *Queue) runJob(job *Job, jp JobParams) {
	atomic.AddInt32(&q.runningjobs, 1)
	go func(targetjob *Job, job JobParams) {
		if job.Delay > 0 {
			targetjob.RunWithDelay(job.Arguments, job.Delay)
		} else if job.Period > 0 {
			targetjob.RunEvery(job.Arguments, job.Period)
		} else {
			targetjob.Run(job.Arguments)
		}

	}(job, jp)
}

func (q *Queue) IsEmpty() bool {
	if len(q.jobs) > 0 {
		return false
	} else {
		return true
	}
}

//Process provides start to processing jobs
func (q *Queue) Process() {
	go func() {
		for {
			for _, job := range q.jobs {
				if job.IsDone() {
					idx := 0
					for i, pname := range q.jobs {
						if pname.Title == job.Title {
							idx = i
							break
						}
					}
					if len(q.jobs) > 0 {
						q.jobs = append(q.jobs[:idx], q.jobs[idx+1:]...)
					}
					info := Info{
						Title:  job.Title,
						Status: 1,
					}

					info.storeInfo(q.dbstore)

					resultitem, err := job.getResult()
					result := Result{
						Title: job.Title,
						Date:  time.Now(),
					}
					//Serialize result, in the case if task contain result value
					if err == nil {
						res, err := json.Marshal(resultitem)
						if err != nil {
							log.Fatal(fmt.Sprintf("Can't get checksum from resut of %s", job.Title))
						}
						result.Result = resultitem
						result.DataChecksum = getChecksum(res)
					}

					result.storeResult(q.dbstore)
					info = Info{
						Title:  job.Title,
						Status: 0,
					}

					info.storeInfo(q.dbstore)
					atomic.AddInt32(&q.runningjobs, -1)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

//ProcessGroups provides loop for processing jobs with type "Groupjobs"
func (q *Queue) ProcessGroups() {
	go func() {
		for {
			for idx, groupjob := range q.groupjobs {
				go func() {
					var wg sync.WaitGroup
					for i := 0; i < len(groupjob); i++ {
						wg.Add(i)
					}
					for i, jobitem := range groupjob {
						jobitem.Run(jobitem.Arguments)
						if jobitem.IsDone() {
							groupjob = append(groupjob[:i], groupjob[i+1:]...)
							wg.Done()
						}
					}
					wg.Wait()
					fmt.Println("All group jobs was completed")
				}()
				q.groupjobs = append(q.groupjobs[:idx], q.groupjobs[idx+1:]...)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}
