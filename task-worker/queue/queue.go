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

func (q *Queue) Dequeue() Job {
	job := <-q.JobChannel
	return job
}
