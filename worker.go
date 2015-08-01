package samovar

import (
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"time"
)

const (
	QUEUENAME = "samowar1"
)

//Function argument to job
type funcarg func(interface{}) interface{}
type funcargint func(int) interface{}
type funcargsstring func(string) interface{}

type Worker struct {
	queues map[string]*Queue
	jobs   map[string]*Job
	host   string
	port   uint
	stop   bool
	//Backend provides comunications with redis
	Backend  *RedisBackend
	jobqueue []*Job
}

//CreateWorker provides initialization of worker
func createWorker(opt *SamovarOptions) *Worker {
	worker := new(Worker)
	worker.host = opt.Host
	worker.port = opt.Port
	worker.queues = map[string]*Queue{}
	worker.Backend = InitRedisBackend()
	if len(opt.Queues) != 0 {
		for _, qname := range opt.Queues {
			worker.AddQueue(qname)
		}
	}

	if !opt.NotDefaultQueue {
		log.Println("initialize default queue")
		worker.AddQueue("default")
	}

	worker.jobs = map[string]*Job{}
	worker.jobqueue = []*Job{}
	return worker
}

//StartWorker provides start server
func (work *Worker) StartWorker() {
	log.Printf("Start worker:")
	work.start()
	detectExit()
	http.ListenAndServe(":8080", nil)

}

//AddJob provides registration of the new job
func (work *Worker) AddJob(title string, fn funcarg) {
	log.Println(fmt.Sprintf("Register new job: %s", title))
	work.jobs[title] = CreateJob(title, fn)
}

//CheckJob provides checking job by title in store
func (work *Worker) CheckJob(title string) bool {
	_, ok := work.jobs[title]
	return ok
}

//This method provides creation of new job queue
func (work *Worker) AddQueue(title string) {
	log.Printf(fmt.Sprintf("Create queue %s", title))
	work.registerQueue(title)
	work.Backend.registerQueue(title)
}

func (work *Worker) registerQueue(title string) {
	queuename := "samovar_" + title
	_, ok := work.queues[queuename]
	if ok {
		log.Printf(fmt.Sprintf("Queue with title %s already exist", title))
	} else {
		work.queues[queuename] = CreateQueue(queuename)
		work.queues[queuename].Process()
		msg := fmt.Sprintf("Queue %s was created", title)
		work.Backend.subscribe(title)
		log.Print(msg)
	}
}

//Before start of worker, show list of jobs
func (work *Worker) showJobs() {
	for _, value := range work.queues {
		fmt.Println(value)
	}
}

func (worker *Worker) start() {
	backend := worker.Backend
	backend.subscribe(QUEUENAME)
	go func() {
		for {

			reply, err := backend.receiveMessages()
			//Get next job
			if err != nil {

			}
			//Need to process new job
			switch msg := reply.(type) {
			case *redis.Message:
				jobobject := msg.Payload
				job := getJobArguments(jobobject)
				if !worker.CheckJob(job.Name) {
					log.Printf(fmt.Sprintf("Job %s not found", job.Name))
				} else {
					targetjob := worker.jobs[job.Name]
					worker.RunNewJob(msg.Channel, targetjob, job)
				}

			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (worker *Worker) RunNewJob(queuename string, tj *Job, jp JobParams) {
	queue, ok := worker.queues[queuename]
	if !ok {
		log.Printf(fmt.Sprintf("Error: queue with the name %s is not found", queuename))
	} else {
		queue.Put(tj, jp)
	}
}

func (worker *Worker) Stop() {
	worker.stop = true
}

func (worker *Worker) Reload() {

}
