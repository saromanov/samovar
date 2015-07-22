package samovar

import ()

type Client struct {
	backend *RedisBackend
}

func InitClient() *Client {
	client := new(Client)
	client.backend = InitRedisBackend()
	return client
}

//Send provides sending arguments to the function
func (gro *Client) Send(jobtitle string, args interface{}) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
	}))
}

func (gro *Client) SendWithDelay(jobtitle string, delay uint, args interface{}) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Delay:     delay,
	}))
}

func (gro *Client) SendWithPeriod(jobtitle string, sec uint, args interface{}) {
	gro.backend.publishJob(PrepareParameters(&JobParams{
		Name:      jobtitle,
		Arguments: args,
		Period:    sec,
	}))
}
