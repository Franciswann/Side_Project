package main

import (
	"fmt"
	"sync"

	"github.com/Franciswann/Side_Project/task-worker/producer"
	"github.com/Franciswann/Side_Project/task-worker/queue"
	"github.com/Franciswann/Side_Project/task-worker/worker"
)

var jobWg sync.WaitGroup
var producerWg sync.WaitGroup

func main() {
	// Job Queue
	q := queue.NewQueue(10)
	// Dead Letter Queue
	dlq := queue.NewQueue(10)

	producer := producer.Producer{JobCount: 5}

	producerWg.Add(1)
	go func() {
		defer producerWg.Done()
		producer.Start(q, &jobWg)
	}()

	worker.StartPool(3, q, dlq, &jobWg)

	producerWg.Wait()
	jobWg.Wait()
	close(q.JobChannel)
	close(dlq.JobChannel)

	fmt.Println("--- Dead Letter Queue ---")
	for failedJob := range dlq.JobChannel {
		fmt.Printf("Failed Job: %v ", failedJob)
	}
}
