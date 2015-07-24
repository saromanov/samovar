package samovar

import (
	"gopkg.in/redis.v3"
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
	}))
}

//SendWithDelay provides sending arguments to job with delay
func (gro *Client) SendWithDelay(jobtitle string, delay uint, args interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}))
}

//SendWithPeriod provides starting periodic task execution
func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args interface{}) {
	gro.backend.publishJob(prepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}))
}

func (gro *Client) GetResult(title string) interface{} {
	return gro.client.HGet("samovar", title)
}
