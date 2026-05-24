package queue

type Job struct {
	ID         int
	Payload    string
	Retries    int
	MaxRetries int
}

type Queue struct {
	JobChannel chan Job
}

func NewQueue(size int) *Queue {
	q := Queue{JobChannel: make(chan Job, size)}
	return &q
}

func (q *Queue) Enqueue(job Job) {
	q.JobChannel <- job
}

// Dequeue blocks until a job is available. Use for range q.JobChannel for continuous consumption.
func (q *Queue) Dequeue() Job {
	job := <-q.JobChannel
	return job
}
