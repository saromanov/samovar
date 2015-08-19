package samovar

import
(
	"errors"
	"fmt"
)

//Jobs provides getting jobs and append jobs
type Jobs struct {
	jobs   map[string]*Job
	groupjobs map[string][]*GroupJob
	limit  int
}


//GetJob provides getting job by title
func (j*Jobs) GetJob(title string, job *Job) error {
	j.jobs[title].numberofcalls++
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

//Append Groupjobs
func (j*Jobs) AppendGroupJob(title string, grjobs[]*GroupJob, reply *bool) error {
	j.groupjobs[title] = grjobs
	*reply = true
	return nil
}