package samovar

import ()

//Options for queue
type QueueOptions struct {
}

//This queue provides basic data structure
type Queue struct {
	jobs    []Job
	title   string
	limit   int
	options *QueueOptions
}

func CreateQueue(title string) *Queue {
	queue := new(Queue)
	queue.jobs = []Job{}
	queue.title = title
	//By default, queue is unlimited
	queue.limit = -1
	queue.options = new(QueueOptions)
	return queue
}

//This method provides getting and running job
func (q *Queue) Get() {

}
