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

//SamovarOptions provides options for the worker
type SamovarOptions struct {
	Queues []string
	host   string
	port   uint
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
func (gro *Samovar) CreateWorker(opt *SamovarOptions) *Worker {
	host := opt.host
	port := opt.port
	if opt.host == "" && opt.port == 0 {
		host = "localhost"
		port = 8080
	}
	gro.worker = CreateWorker(host, port)
	return gro.worker
}
