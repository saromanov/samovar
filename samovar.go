package samovar

const (
	//Default queue name
	defqueue = "default"
)

type Samovar struct {
	queues map[string]*Queue
	worker *Worker
	client *Client
}

//Init provides initialization of basic object
func Init() *Samovar {
	gro := new(Samovar)
	return gro
}

//CreateClient provides initialization of client
func (gro *Samovar) CreateClient() *Client {
	gro.client = InitClient()
	return gro.client
}

//CreateWorker provides initialization of the worker
func (gro *Samovar) CreateWorker(host string, port uint) *Worker {
	gro.worker = CreateWorker(host, port)
	return gro.worker
}
