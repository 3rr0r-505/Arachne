package crawler

type Job struct {
	URL   string
	Depth int
}

type JobQ struct {
	jobs chan Job
}

func NewJobQueue(size int) *JobQ {
	return &JobQ{
		jobs: make(chan Job, size),
	}
}

func (jq *JobQ) Push(job Job) {
	jq.jobs <- job
}

func (jq *JobQ) Pop() (Job, bool) {
	job, ok := <-jq.jobs
	return job, ok
}
