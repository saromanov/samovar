package samovar

import (
	"gopkg.in/redis.v3"
	"time"
	"net"
	"net/rpc"
	"log"
	"fmt"
)

type Client struct {
	backend *RedisBackend
	client  *redis.Client
}

type JobOptions struct {
	Title     string
	Arguments interface{}
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
func (gro *Client) SendWithDelay(jobtitle string, delay uint, args interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}), "default")
}

//SendWithPeriod provides starting periodic task execution
func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}), "default")
}

//GetResult provides non-async version if getting results from the job
func (gro *Client) GetResult(title string) interface{} {
	return getResult(gro.client, title)
}

//GetStat provides statistics for the job with title
func (gro *Client) GetJobItem(title string) {
	callAsRPC("Jobs.GetJob", title)
}

func resolveQueueName(title string) string {
	if title == "" {
		return "default"
	}
	return title
}

func callAsRPC(name, title string) {
	timeout := time.Duration(100) * time.Millisecond

	item, err := net.DialTimeout("tcp", ADDR, timeout)
	if err != nil {
		log.Fatal(err)
	}

	connection := rpc.NewClient(item)
	var job Job
	errcall := connection.Call(name, title, &job)
	if errcall != nil {
		log.Fatal(errcall)
	}

	fmt.Println(job)
}
