package worker

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

type Worker struct {
	ID int
}

func (w *Worker) Start(q *queue.Queue, dlq *queue.Queue, wg *sync.WaitGroup) {
	for task := range q.JobChannel {
		// 30% fail
		if rand.IntN(10) < 3 {
			if task.Retries >= task.MaxRetries {
				dlq.Enqueue(task)
				wg.Done()
				continue
			}
			// retry
			task.Retries++
			q.Enqueue(task)
			continue
		}
		// process successfully
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Worker %v processing: %v\n", w.ID, task)
		wg.Done()
	}
}

// Worker Pool
func StartPool(numWorkers int, q *queue.Queue, dlq *queue.Queue, wg *sync.WaitGroup) {
	for number := 0; number < numWorkers; number++ {
		worker := Worker{ID: number + 1}

		go func(w Worker) {
			w.Start(q, dlq, wg)
		}(worker)
	}
}
