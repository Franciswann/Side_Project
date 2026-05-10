package queue

type Job struct {
	ID      int
	Payload string
}

type Queue struct {
	JobChannel chan Job
}

func NewQueue(size int) *Queue {
	q := Queue{JobChannel: make(chan Job, size)}
	return &q
}
