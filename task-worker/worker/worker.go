package worker

import (
	"fmt"
	"sync"

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

// Worker Pool
func StartPool(numWorkers int, q *queue.Queue, wg *sync.WaitGroup) {
	wg.Add(numWorkers)
	for number := range numWorkers {
		worker := Worker{ID: number + 1}
		go func() {
			defer wg.Done()
			worker.Start(q)
		}()
	}
}
