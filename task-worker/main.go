package main

import (
	"fmt"

	"github.com/Franciswann/Side_Project/task-worker/queue"
)

func main() {
	q := queue.NewQueue(10)
	task1 := queue.Job{ID: 1, Payload: "task-1"}
	task2 := queue.Job{ID: 2, Payload: "task-2"}
	task3 := queue.Job{ID: 3, Payload: "task-3"}
	q.Enqueue(task1)
	q.Enqueue(task2)
	q.Enqueue(task3)
	close(q.JobChannel)

	for task := range q.JobChannel {
		fmt.Printf("%v\n", task)
	}

}
