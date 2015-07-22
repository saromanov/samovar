package samovar

import ()

type Client struct {
	backend *RedisBackend
}

type JobOptions struct {
	Arguments interface{}
	Delay     uint
	Period    uint
	Priority  uint
}

//Init client provides initialization of samovar client
func InitClient() *Client {
	client := new(Client)
	client.backend = InitRedisBackend()
	return client
}

//Send provides sending arguments to the function
func (gro *Client) Send(jobtitle string, opt *JobOptions) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: opt.Arguments,
	}))
}

//SendWithDelay provides sending arguments to job with delay
func (gro *Client) SendWithDelay(jobtitle string, delay uint, args interface{}) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}))
}

//SendWithPeriod provides starting periodic task execution
func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args interface{}) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}))
}
