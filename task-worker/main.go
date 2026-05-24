package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Franciswann/Side_Project/task-worker/producer"
	"github.com/Franciswann/Side_Project/task-worker/queue"
	"github.com/Franciswann/Side_Project/task-worker/worker"
)

func main() {
	var jobWg sync.WaitGroup
	var producerWg sync.WaitGroup

	q := queue.NewQueue(10)
	dlq := queue.NewQueue(10)

	// context is used to gracefully stop the producer on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	producer := producer.Producer{JobCount: 5}

	producerWg.Add(1)

	go func() {
		defer producerWg.Done()
		producer.Start(ctx, q, &jobWg)
	}()

	worker.StartPool(3, q, dlq, &jobWg)

	// Listen for OS shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Gracefully stop the producer when shutdown signal is received
	go func() {
		defer close(sigChan)
		<-sigChan
		cancel()
		fmt.Println("\nShutdown signal received, waiting for jobs to finish...")
	}()

	producerWg.Wait()
	jobWg.Wait()

	close(q.JobChannel)
	close(dlq.JobChannel)

	fmt.Println("--- Dead Letter Queue ---")
	for failedJob := range dlq.JobChannel {
		fmt.Printf("Failed Job: %v ", failedJob)
	}
}
