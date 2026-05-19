package main

import (
	"sync"

	"github.com/Franciswann/Side_Project/task-worker/queue"
	"github.com/Franciswann/Side_Project/task-worker/worker"
)

var wg sync.WaitGroup

func main() {

	q := queue.NewQueue(10)
	task1 := queue.Job{ID: 1, Payload: "task-1"}
	task2 := queue.Job{ID: 2, Payload: "task-2"}
	task3 := queue.Job{ID: 3, Payload: "task-3"}
	q.Enqueue(task1)
	q.Enqueue(task2)
	q.Enqueue(task3)
	close(q.JobChannel)

	worker := worker.Worker{ID: 1}
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.Start(q)
	}()
	wg.Wait()
}
