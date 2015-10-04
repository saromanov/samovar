package backend

type Backend interface {

	//Registration of the new queue
	RegisterQueue(queuename string)

	//Receive new message
	ReceiveMessages() (interface{}, error)

	//Subscribe to the new queue
	Subscribe(queuename string)

	//Push new job
	PublishJob(jobtitle string, queuename string)
}