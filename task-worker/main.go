package main

import (
	"fmt"
	"sync"

	"github.com/Franciswann/Side_Project/task-worker/queue"
	"github.com/Franciswann/Side_Project/task-worker/worker"
)

var wg sync.WaitGroup

func main() {

	q := queue.NewQueue(10)
	task1 := queue.Job{ID: 1, Payload: "task-1", MaxRetries: 3}
	task2 := queue.Job{ID: 2, Payload: "task-2", MaxRetries: 3}
	task3 := queue.Job{ID: 3, Payload: "task-3", MaxRetries: 3}
	task4 := queue.Job{ID: 4, Payload: "task-4", Retries: 2, MaxRetries: 3}
	task5 := queue.Job{ID: 5, Payload: "task-5", MaxRetries: 3}
	q.Enqueue(task1)
	q.Enqueue(task2)
	q.Enqueue(task3)
	q.Enqueue(task4)
	q.Enqueue(task5)
	// add jobs into JobChannel
	wg.Add(len(q.JobChannel))

	// Dead Letter Queue
	dlq := queue.NewQueue(10)

	worker.StartPool(3, q, dlq, &wg)

	wg.Wait()
	close(q.JobChannel)
	close(dlq.JobChannel)

	fmt.Println("--- Dead Letter Queue ---")
	for failedJob := range dlq.JobChannel {
		fmt.Printf("Failed Job: %v ", failedJob)
	}
}
