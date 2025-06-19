// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// // Job is the basic unit of work. For now, it's just a function.
// type Job struct {
// 	Task func()
// }

// // Scheduler manages a pool of workers to execute jobs.
// type Scheduler struct {
// 	jobQueue    chan Job
// 	workerCount int
// 	wg          sync.WaitGroup
// }

// // NewScheduler creates a new scheduler and starts the worker pool.
// func NewScheduler(workerCount int) *Scheduler {
// 	s := &Scheduler{
// 		// A buffered channel allows us to add jobs without waiting for a worker to be free.
// 		jobQueue:    make(chan Job, 100),
// 		workerCount: workerCount,
// 	}

// 	// Start the workers
// 	s.wg.Add(s.workerCount)
// 	for i := 0; i < s.workerCount; i++ {
// 		go s.worker(i + 1)
// 	}

// 	fmt.Printf("Scheduler started with %d workers.\n", s.workerCount)
// 	return s
// }

// // worker is a goroutine that continuously processes jobs from the jobQueue.
// func (s *Scheduler) worker(id int) {
// 	defer s.wg.Done()
// 	fmt.Printf("Worker %d started\n", id)
// 	// The for-range loop on a channel will automatically break when the channel is closed.
// 	for job := range s.jobQueue {
// 		fmt.Printf("Worker %d is starting job\n", id)
// 		job.Task()
// 		fmt.Printf("Worker %d finished job\n", id)
// 	}
// 	fmt.Printf("Worker %d shutting down\n", id)
// }

// // AddJob adds a new job to the scheduler's queue.
// func (s *Scheduler) AddJob(task func()) {
// 	job := Job{Task: task}
// 	s.jobQueue <- job
// }

// // Stop gracefully shuts down the scheduler and its workers.
// func (s *Scheduler) Stop() {
// 	// Closing the channel signals the workers to stop processing new jobs.
// 	close(s.jobQueue)
// 	// Wait for all workers to finish their current job and exit.
// 	s.wg.Wait()
// 	fmt.Println("Scheduler stopped.")
// }

// func main() {
// 	// Create a scheduler with 3 concurrent workers.
// 	scheduler := NewScheduler(3)

// 	// Add some jobs.
// 	for i := 1; i <= 5; i++ {
// 		// We need to capture the loop variable `i` in the closure.
// 		jobID := i
// 		scheduler.AddJob(func() {
// 			fmt.Printf("Executing job %d...\n", jobID)
// 			time.Sleep(2 * time.Second) // Simulate work
// 		})
// 	}

// 	// Give some time for jobs to be processed before stopping.
// 	// In a real app, the scheduler would run for the lifetime of the application.
// 	time.Sleep(10 * time.Second)

// 	scheduler.Stop()
// }
