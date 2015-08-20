package samovar

import (
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	//"net/http"
	"time"
)

const (
	QUEUENAME = "samowar1"
)

type Worker struct {
	queues map[string]*Queue
	//jobs   map[string]*Job
	host string
	port uint
	stop bool
	//Backend provides comunications with redis
	Backend  *RedisBackend
	jobqueue []*Job
	jobs     *Jobs
}

type GroupJob struct {
	Title string
	Fn    interface{}
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

	worker.jobs = &Jobs{
		jobs: map[string]*Job{},
		groupjobs: map[string][]*Job{},
	}
	worker.jobqueue = []*Job{}
	return worker
}

//StartWorker provides start server
func (work *Worker) StartWorker() {
	log.Printf("Start worker:")
	work.start()
	detectExit()
	//InitRPC("").Run()
	StartServer(work.jobs)

}

//AddJob provides registration of the new job
func (work *Worker) AddJob(title string, fn interface{}) {
	log.Println(fmt.Sprintf("Register new job: %s", title))
	if work.jobs.Exists(title) {
		log.Fatal(fmt.Sprintf("Job with the title %s already registred", title))
	}
	var reply bool
	err := work.jobs.AppendJob(title, CreateJob(title, fn), &reply)
	if err != nil {
		log.Fatal(err)
	}
}

//AddGroupJobs provides append group of depended of each other jobs
func (work *Worker) AddGroupJobs(title string, groupjobs []*GroupJob) {
	log.Println(fmt.Sprintf("Register new group of jobs: %s", title))
	if work.jobs.Exists(title) {
		log.Fatal(fmt.Sprintf("Job with the title %s already registred", title))
	}
	var reply bool
	err := work.jobs.AppendGroupJob(title, groupjobs, &reply)
	if err != nil {
		log.Fatal(err)
	}

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
		work.queues[queuename].ProcessGroups()
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
				var targetjob []*Job
				//This method provides getting job or list > 1, getting group of jobs
				err := worker.jobs.GetJob(job.Name, &targetjob)
				if err != nil {
					log.Fatal(err)
				}

				if len(targetjob) == 1 {
					worker.jobs.increment(job.Name)
					worker.RunNewJob(msg.Channel, targetjob[0], job)
				} else {
					worker.RunNewJobGroup(msg.Channel, targetjob)
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

func (worker *Worker) RunNewJobGroup(queuename string, tj[]*Job) {
	queue, ok := worker.queues[queuename]
	if !ok {
		log.Printf(fmt.Sprintf("Error: queue with the name %s is not found", queuename))
		} else {
			queue.PutGroup(tj)
		}
}

func (worker *Worker) Stop() {
	worker.stop = true
}

func (worker *Worker) Reload() {

}
