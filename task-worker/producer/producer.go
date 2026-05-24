package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

type Producer struct {
	JobCount int
}

func (p *Producer) Start(ctx context.Context, q *queue.Queue, wg *sync.WaitGroup) {
	for job := 1; job <= p.JobCount; job++ {
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(1)
			task := queue.Job{ID: job, Payload: fmt.Sprintf("task-%v", job), Retries: 0, MaxRetries: 3}
			fmt.Printf("Producer sending: %v\n", task)
			time.Sleep(300 * time.Millisecond)
			q.Enqueue(task)
		}
	}
}
