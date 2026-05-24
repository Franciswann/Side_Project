package worker

import (
	"sync"
	"testing"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

var wg sync.WaitGroup

func TestStart(t *testing.T) {
	q := queue.NewQueue(10)
	dlq := queue.NewQueue(10)

	job1 := queue.Job{ID: 1, Payload: "task-1", Retries: 0, MaxRetries: 3}
	job2 := queue.Job{ID: 2, Payload: "task-2", Retries: 3, MaxRetries: 3}

	wg.Add(1)

	worker := Worker{
		ShouldFail: func() bool {
			return true
		},
	}

	worker.ProcessJob(job1, q, dlq, &wg)
	worker.ProcessJob(job2, q, dlq, &wg)

	close(q.JobChannel)
	close(dlq.JobChannel)

	t.Run("failed job with retries remaining is re-enqueued with incremented retry count", func(t *testing.T) {
		got := <-q.JobChannel
		want := queue.Job{ID: 1, Payload: "task-1", Retries: 1, MaxRetries: 3}
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("failed job with no retries remaining is sent to dead letter queue", func(t *testing.T) {
		got := <-dlq.JobChannel
		want := job2
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
