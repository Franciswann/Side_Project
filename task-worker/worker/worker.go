package worker

import (
	"fmt"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

type Worker struct {
	ID int
}

func (w *Worker) Start(q *queue.Queue) {
	for task := range q.JobChannel {
		fmt.Printf("Worker %v processing: %v\n", w.ID, task)
	}
}
