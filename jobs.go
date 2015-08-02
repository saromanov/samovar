package samovar

import
(
	"errors"
	"fmt"
)

//Jobs provides getting jobs and append jobs
type Jobs struct {
	jobs   map[string]*Job
	limit  int
}


func (j*Jobs) GetJob(title string, job *Job) error {
	currentjob, ok := j.jobs[title]
	if !ok {
		return errors.New(fmt.Sprintf("Job with the title %s is not found", title))
	}
	if job == nil {
		return nil
	}
	*job = *currentjob
	return nil
}

//AppendJob provides set new job
func (j*Jobs) AppendJob(title string, job*Job, reply *bool) error {
	if j.limit > 0 && len(j.jobs) == j.limit {
		*reply = false
		return errors.New("Can't add new job because of limit of jobs")
	}

	j.jobs[title] = job
	*reply = true
	return nil
}