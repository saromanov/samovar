package samovar

import (
	"time"
	"gopkg.in/redis.v3"
	"encoding/json"
	"log"
	"fmt"
)

//Options for queue
type QueueOptions struct {

}

//This queue provides basic data structure
type Queue struct {
	jobs    []*Job
	title   string
	limit   int
	options *QueueOptions
	dbstore  *redis.Client
}

func CreateQueue(title string) *Queue {
	queue := new(Queue)
	queue.jobs = []*Job{}
	queue.title = title
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

func (q *Queue) runJob(job *Job, jp JobParams) {
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

func (q *Queue) find(title string) {

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
					info := Info {
						Title: job.Title,
						Status: 1,
					}

					info.storeInfo(q.dbstore)

					resultitem := job.getResult()
					res, err := json.Marshal(resultitem)
					if err != nil {
						log.Fatal(fmt.Sprintf("Can't get checksum from resut of %s", job.Title))
					}
					result := Result{
						Title:  job.Title,
						Result: resultitem,
						Date: time.Now(),
						DataChecksum: getChecksum(res),
					}
					result.storeResult(q.dbstore)
					info = Info {
						Title: job.Title,
						Status: 0,
					}

					info.storeInfo(q.dbstore)
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}
