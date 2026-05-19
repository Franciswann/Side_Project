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
	task4 := queue.Job{ID: 4, Payload: "task-4"}
	task5 := queue.Job{ID: 5, Payload: "task-5"}
	q.Enqueue(task1)
	q.Enqueue(task2)
	q.Enqueue(task3)
	q.Enqueue(task4)
	q.Enqueue(task5)
	close(q.JobChannel)

	worker.StartPool(3, q, &wg)
	wg.Wait()
}
