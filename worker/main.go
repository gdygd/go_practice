package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"worker_test/worker"
)

type UserJob struct {
	NumId int
}

// implement worker.Job interface
func (j *UserJob) Run(ctx context.Context) error {
	fmt.Printf("#1 Job number : %d \n", j.NumId)
	time.Sleep(time.Second * 5)
	fmt.Printf("#2 End Job number : %d \n", j.NumId)

	return nil
}

func MakeJobs() []UserJob {
	j := []UserJob{}
	for i := 0; i < 10; i++ {
		j = append(j, UserJob{i + 1})
	}
	return j
}

func main() {
	wp := worker.NewPool(10, 10)
	wp.Start()

	jobs := MakeJobs()
	for _, j := range jobs {
		// if i%2 == 0 {
		// 	wp.Submit(&j)
		// } else {
		// 	wp.SubmitWithTimeout(&j, 3)
		// }
		wp.SubmitWithTimeout(&j, 3)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wp.Stop(ctx); err != nil {
		log.Printf("stop error : %v", err)
	}
}
