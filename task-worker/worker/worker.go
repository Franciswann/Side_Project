package worker

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

type Worker struct {
	ID int
}

func (w *Worker) Start(q *queue.Queue, wg *sync.WaitGroup) {
	for task := range q.JobChannel {
		// 30% fail
		if rand.IntN(10) < 3 {
			if task.Retries >= task.MaxRetries {
				fmt.Printf("Job {%v} failed permanently\n", task.ID)
				wg.Done()
				continue
			}
			// retry
			task.Retries++
			q.Enqueue(task)
			continue
		}
		// process successfully
		fmt.Printf("Worker %v processing: %v\n", w.ID, task)
		wg.Done()
	}
}

// Worker Pool
func StartPool(numWorkers int, q *queue.Queue, wg *sync.WaitGroup) {
	for number := 0; number < numWorkers; number++ {
		worker := Worker{ID: number + 1}

		go func(w Worker) {
			w.Start(q, wg)
		}(worker)
	}
}
