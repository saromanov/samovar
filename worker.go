package samovar

import (
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"time"
)

//Function argument to job
type funcarg func(interface{}) interface{}
type funcargint func(int) interface{}
type funcargsstring func(string) interface{}

type Worker struct {
	queues  map[string]*Queue
	jobs    map[string]*Job
	host    string
	port    uint
	stop    bool
	Backend *RedisBackend
}

func CreateWorker(host string, port uint) *Worker {
	worker := new(Worker)
	worker.host = host
	worker.port = port
	worker.queues = map[string]*Queue{}
	log.Println("initialize default queue")
	worker.Backend = InitRedisBackend()
	worker.AddQueue("default")
	worker.jobs = map[string]*Job{}
	return worker
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (work *Worker) StartWorker() {
	log.Printf("Start worker:")
	work.start()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

//AddJob provides registration of the new job
func (work *Worker) AddJob(title string, fn funcarg) {
	log.Println("Create new job")
	work.jobs[title] = CreateJob(title, fn)
}

//CheckJob provides checking job by title in store
func (work *Worker) CheckJob(title string) bool {
	_, ok := work.jobs[title]
	return ok
}

//This method provides creation of new job queue
func (work *Worker) AddQueue(title string) {
	work.queues[title] = CreateQueue(title)
	work.Backend.registerQueue(title)
}

func (work *Worker) registerQueue(title string) {
	//Нужно оставить только subscrube
	_, ok := work.queues[title]
	if ok {
		log.Printf(fmt.Sprintf("Queue with title %s already exist", title))
	} else {
		work.queues[title] = CreateQueue(title)
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
	backend.subscribe("samowar1")
	go func() {
		for {

			reply, err := backend.ReceiveMessages()
			//Get next job
			if err != nil {

			}
			//Need to process new job
			switch msg := reply.(type) {
			case *redis.Message:
				jobobject := msg.Payload
				job := GetJobArguments(jobobject)
				if !worker.CheckJob(job.Name) {
					log.Printf(fmt.Sprintf("Job %s not found", job.Name))
				} else {
					targetjob := worker.jobs[job.Name]
					if job.Delay > 0 {
						targetjob.RunWithDelay(job.Arguments, job.Delay)
					} else if job.Period > 0 {
						targetjob.RunEvery(job.Arguments, job.Period)
					} else {
						targetjob.Run(job.Arguments)
					}
				}

			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (worker *Worker) Stop() {
	worker.stop = true
}

func (worker *Worker) Reload() {

}
