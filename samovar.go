package samovar

type Samovar struct {
	queues map[string]*Queue
	worker *Worker
	client *Client
}

//SamovarOptions provides options for the worker
type SamovarOptions struct {
	Queues       []string
	Host         string
	Port         uint
	NotDefaultQueue bool
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
	if opt.Host == "" && opt.Port == 0 {
		opt.Host = "localhost"
		opt.Port = 8080
	}
	gro.worker = createWorker(opt)
	return gro.worker
}
