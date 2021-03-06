package samovar

import (
	"./backend"
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type Client struct {
	back      backend.Backend
	client    *redis.Client
	rpcclient *rpc.Client
	store     *Store
}

type JobOptions struct {
	Title     string
	Arguments interface{}
	Delay     uint
	Period    uint
	Priority  uint
	Queue     string
}

//JobItem provides information about job, getting from server
type ResultJob struct {
	NumberOfCalls int
	Done          bool
	Result        interface{}
	ExecutionTime float64
}

type AsyncResult struct {
	Result interface{}
}

//Init client provides initialization of samovar client
func InitClient() *Client {
	client := new(Client)
	client.back = backend.InitRedisBackend()
	client.client = backend.InitRedis("localhost:6379")
	timeout := time.Duration(100) * time.Millisecond
	item, err := net.DialTimeout("tcp", ADDR, timeout)
	if err != nil {
		log.Fatal(err)
	}

	client.rpcclient = rpc.NewClient(item)
	client.store = InitStore()
	return client
}

//Send provides sending arguments to the function
func (gro *Client) Send(jobtitle string, opt *JobOptions) string {
	id := RandString(jobidsize)
	gro.back.PublishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: opt.Arguments,
		JobID:     id,
	}), resolveQueueName(opt.Queue))

	return fmt.Sprintf("Job ID: %s", id)
}

func (gro *Client) SendAsync(jobtitle string, id string, opt *JobOptions) *AsyncResult {
	gro.Send(jobtitle, opt)

	//Wait until job will be complete
	for {
		info, _ := getInfo(gro.client, jobtitle)
		if info.Status == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return &AsyncResult{Result: gro.GetResult(jobtitle, id)}
}

//SendWithDelay provides sending arguments to job with delay
func (gro *Client) SendWithDelay(jobtitle string, delay uint, args []interface{}) {
	gro.back.PublishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}), "default")
}

//SendWithPeriod provides starting periodic task execution
func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args []interface{}) {
	gro.back.PublishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}), "default")
}

//SendMany provides starting several jobs
//which have been registered as single jobs
func (gro *Client) SendMany(jobs []*JobOptions) {
	if len(jobs) == 0 {
		log.Printf("Number of sending tasks is zero")
		return
	}

	//Get tasks randomly
	for {
		if len(jobs) == 0 {
			break
		}

		idx := rand.Intn(len(jobs))
		job := jobs[idx]
		jobs = append(jobs[:idx], jobs[idx+1:]...)
		queuename := "default"
		if job.Queue != "" {
			queuename = job.Queue
		}
		gro.back.PublishJob(prepareParameters(&JobParams{
			Name:      job.Title,
			Arguments: job.Arguments,
		}), queuename)
	}

}

//GetResult provides non-async version if getting results from the job
func (gro *Client) GetResult(title string, id string) interface{} {
	result, err := getResult(gro.client, fmt.Sprintf("%s%s", title, id))
	if err != nil {
		log.Fatal(err)
	}
	return result
}

//GetResultByID provides getting result from past job by ID
func (gro *Client) GetResultByID(id string) interface{} {
	result, err := getResultById(gro.client, id)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

//Saveresult provides storing result of task by title "tasktile" to key-value store
//Note: Now suppoted only redis
func (gro *Client) SaveResult(tasktitle, id, key string) {
	result := gro.GetResult(tasktitle, id)
	gro.store.Set(key, string(Marshal(result)))

}

//GetStat provides statistics for the job with title
func (gro *Client) GetJobItem(title string) *ResultJob {
	var numcals int
	errcall := gro.rpcclient.Call("Jobs.GetNumberOfCalls", title, &numcals)
	if errcall != nil {
		log.Fatal(errcall)
	}
	var done bool
	errcall2 := gro.rpcclient.Call("Jobs.IsJobDone", title, &done)
	if errcall2 != nil {
		log.Fatal(errcall2)
	}
	var execution float64
	errcall3 := gro.rpcclient.Call("Jobs.ExecutionTime", title, &execution)
	if errcall3 != nil {
		log.Fatal(errcall3)
	}
	return &ResultJob{NumberOfCalls: numcals, Done: done, ExecutionTime: execution}
}

func resolveQueueName(title string) string {
	if title == "" {
		return "default"
	}
	return title
}
