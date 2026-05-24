package worker

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

type Worker struct {
	ID         int
	ShouldFail func() bool
}

func (w *Worker) Start(q *queue.Queue, dlq *queue.Queue, wg *sync.WaitGroup) {
	for task := range q.JobChannel {
		w.ProcessJob(task, q, dlq, wg)
	}
}

func (w *Worker) ProcessJob(task queue.Job, q *queue.Queue, dlq *queue.Queue, wg *sync.WaitGroup) {
	if w.ShouldFail() {
		if task.Retries >= task.MaxRetries {
			dlq.Enqueue(task)
			wg.Done()
			return
		}
		// re-enqueue for retry, wg count remains the same
		task.Retries++
		q.Enqueue(task)
		return
	}
	// process successfully
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Worker %v processing: %v\n", w.ID, task)
	wg.Done()
}

// Worker Pool
func StartPool(numWorkers int, q *queue.Queue, dlq *queue.Queue, wg *sync.WaitGroup) {
	for number := 0; number < numWorkers; number++ {
		worker := Worker{
			ID: number + 1,
			ShouldFail: func() bool {
				// simulates 30% job failure rate
				return rand.IntN(10) < 3
			},
		}

		go func(w Worker) {
			w.Start(q, dlq, wg)
		}(worker)
	}
}
