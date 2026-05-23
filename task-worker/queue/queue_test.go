package queue

import (
	"reflect"
	"testing"
)

func TestEnqueueDequeue(t *testing.T) {

	t.Run("enqueue and dequeue returns the same job", func(t *testing.T) {
		q := NewQueue(10)
		task1 := Job{ID: 1, Payload: "task-1"}
		q.Enqueue(task1)

		got := q.Dequeue()

		want := Job{ID: 1, Payload: "task-1"}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("multiple jobs are dequeued in FIFO order", func(t *testing.T) {
		q := NewQueue(10)
		task1 := Job{ID: 1}
		task2 := Job{ID: 2}
		task3 := Job{ID: 3}
		q.Enqueue(task1)
		q.Enqueue(task2)
		q.Enqueue(task3)
		close(q.JobChannel)

		jobID := []int{}
		for id := range q.JobChannel {
			jobID = append(jobID, id.ID)
		}

		got := jobID
		want := []int{1, 2, 3}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
