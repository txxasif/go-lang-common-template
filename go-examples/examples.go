// Worker Pool Pattern Visualization:
//
// Jobs Channel (buffered size 5)      Workers (3 goroutines)     Results Channel (buffered size 5)
// ┌──────────────────────┐           ┌─────────────┐            ┌──────────────────────┐
// │ Job{1}  Job{2} ...   │─────┬────►│  Worker 1   │──┐        │ Result{1,2} ...      │
// └──────────────────────┘     │     └─────────────┘  │        └──────────────────────┘
//                              │     ┌─────────────┐   ├───────►
//                              ├────►│  Worker 2   │───┤
//                              │     └─────────────┘   │
//                              │     ┌─────────────┐   │
//                              └────►│  Worker 3   │───┘
//                                    └─────────────┘
//
// Flow:
// 1. Main goroutine sends 5 jobs to jobs channel
// 2. 3 workers process jobs concurrently
// 3. Results collected in results channel
// 4. Main goroutine reads and prints results

package main

import (
	"context"
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Value int
}

func workerPool(ctx context.Context, jobs <-chan Job, results chan<- Result) {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					results <- Result{JobID: job.ID, Value: job.ID * 2}
				}
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 5)
	results := make(chan Result, 5)

	go workerPool(ctx, jobs, results)

	for i := 1; i <= 5; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs)

	for r := range results {
		fmt.Printf("Job %d: %d\n", r.JobID, r.Value)
	}
}
