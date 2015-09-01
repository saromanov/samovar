package samovar

import (
	"gopkg.in/redis.v3"
	"time"
	"net"
	"net/rpc"
	"log"
	"math/rand"
)

type Client struct {
	backend *RedisBackend
	client  *redis.Client
	rpcclient *rpc.Client
}

type JobOptions struct {
	Title     string
	Arguments []interface{}
	Delay     uint
	Period    uint
	Priority  uint
	Queue     string
}

//Init client provides initialization of samovar client
func InitClient() *Client {
	client := new(Client)
	client.backend = InitRedisBackend()
	client.client = initRedis("localhost:6379")
	timeout := time.Duration(100) * time.Millisecond
	item, err := net.DialTimeout("tcp", ADDR, timeout)
	if err != nil {
		log.Fatal(err)
	}

	client.rpcclient = rpc.NewClient(item)
	return client
}

//Send provides sending arguments to the function
func (gro *Client) Send(jobtitle string, opt *JobOptions) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: opt.Arguments,
	}), resolveQueueName(opt.Queue))
}

//SendWithDelay provides sending arguments to job with delay
func (gro *Client) SendWithDelay(jobtitle string, delay uint, args []interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}), "default")
}

//SendWithPeriod provides starting periodic task execution
func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args []interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}), "default")
}


//SendMany provides starting several jobs 
//which have been registered as single jobs
func (gro *Client) SendMany(jobs[] *JobOptions) {
	if len(jobs) == 0 {
		log.Printf("Number of sending tasks is zero")
		return
	}

	//Get tasks randomly
	for  {
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
		gro.backend.publishJob(prepareParameters(&JobParams {
			Name: job.Title,
			Arguments: job.Arguments,
		}), queuename)
	}

}

//GetResult provides non-async version if getting results from the job
func (gro *Client) GetResult(title string) interface{} {
	return getResult(gro.client, title)
}

//GetStat provides statistics for the job with title
func (gro *Client) GetJobItem(title string) {
	var numcals int
	errcall := gro.rpcclient.Call("Jobs.GetNumberOfCalls", title, &numcals)
	if errcall != nil {
		log.Fatal(errcall)
	}
}

func resolveQueueName(title string) string {
	if title == "" {
		return "default"
	}
	return title
}

