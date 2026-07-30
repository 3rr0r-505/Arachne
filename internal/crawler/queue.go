package crawler

type Job struct {
	URL   string
	Depth int
}

type JobQ struct {
	in  chan Job
	out chan Job
}

func NewJobQueue() *JobQ {
	jq := &JobQ{
		in:  make(chan Job),
		out: make(chan Job),
	}
	go jq.run()
	return jq
}

func (jq *JobQ) run() {
	var buffer []Job
	in := jq.in // local copy so we can nil it out later without affecting jq.in

	for {
		if len(buffer) == 0 {
			if in == nil {
				// input closed AND buffer drained — fully done
				close(jq.out)
				return
			}
			job, ok := <-in
			if !ok {
				in = nil // stop trying to receive from a closed channel
				continue
			}
			buffer = append(buffer, job)
		} else {
			select {
			case job, ok := <-in:
				if !ok {
					in = nil
					continue
				}
				buffer = append(buffer, job)
			case jq.out <- buffer[0]:
				buffer = buffer[1:]
			}
		}
	}
}

func (jq *JobQ) Push(job Job) {
	jq.in <- job
}

func (jq *JobQ) Pop() (Job, bool) {
	job, ok := <-jq.out
	return job, ok
}

func (jq *JobQ) Close() {
	close(jq.in)
}
