package samovar

import
(
	"errors"
	"fmt"
)

//Jobs provides getting jobs and append jobs
type Jobs struct {
	jobs   map[string]*Job
	groupjobs map[string][]*Job
	limit  int
}


//GetJob provides getting job by title
func (j*Jobs) GetJob(title string, job*[]*Job) error {
	currentjob, ok := j.jobs[title]
	grjob, ok2 := j.groupjobs[title]
	if !ok && !ok2 {
		return errors.New(fmt.Sprintf("Job with the title %s is not found", title))
	}
	if job == nil {
		return nil
	}
	if ok {
		jobs1 := []*Job{currentjob}
		*job = jobs1
	}

	if ok2 {
		*job = grjob
	}
	return nil
}

func (j*Jobs) GetNumberOfCalls(title string, num *int) error{
	var current int
	currentjob, ok := j.jobs[title]
	if !ok {
		current = 0
	} else {
		current = currentjob.numberofcalls
	}
	*num = current
	return nil
}

//IsJobDone returns if current job is done
func (j *Jobs) IsJobDone(title string, result *bool) error {
	var current bool
	currentjob, ok := j.jobs[title]
	if !ok {
		current = false
	} else {
		current = currentjob.done
	}
	*result = current
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
	newjobs := []*Job{}
	for _, item := range grjobs {
		newjobs = append(newjobs, CreateJob(item.Title, item.Fn))
	}
	j.groupjobs[title] = newjobs
	*reply = true
	return nil
}


//Exists return true of jobs or groupjobs already register
//and false otherwise
func (j*Jobs) Exists(title string) bool{
	_, ok := j.jobs[title]
	if ok {
		return true
	}

	_, ok2 := j.groupjobs[title]
	if ok2 {
		return true
	}

	return false
}

//MeanET returns mean execution time of the job
func (j *Jobs) MeanET(title string)float64 {
	job, ok := j.jobs[title]
	if !ok {
		return 0
	}
	var meanvalue float64
	meanvalue = 0.0
	count := len(job.executionTimes)
	if count == 0 {
		return 0
	}

	for _, value := range job.executionTimes {
		meanvalue += value
	}

	return meanvalue/float64(count)
}


//SHortestJob returns job with shortest execution time
func (j *Jobs) ShortestJob()*Job {
	if len(j.jobs) == 0 {
		return &Job{}
	}

	var value float64
	var job *Job
	value = 999999
	for key, _ := range j.jobs {
		mv := j.MeanET(key)
		if mv < value {
			value = mv
			jobtmp, _ := j.jobs[key]
			job = jobtmp
		}
	}
	return job
}