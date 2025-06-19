package main

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	Task     func()
	Priority int
}
type JobPriorityQueue []*Job

func (pq JobPriorityQueue) Len() int           { return len(pq) }
func (pq JobPriorityQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }

func (pq JobPriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *JobPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Job))
}

func (pq *JobPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	(*pq)[n-1] = nil
	*pq = old[:n-1]
	return item
}

// Scheduler manages prioritized jobs and a worker pool.
type Scheduler struct {
	jobQueue      chan *Job
	priorityQueue JobPriorityQueue
	mu            sync.Mutex
	workerCount   int
	wg            sync.WaitGroup
	stopChan      chan struct{}
	paused        atomic.Bool
}

// NewScheduler creates and starts the scheduler.
func NewScheduler(workerCount int) *Scheduler {
	s := &Scheduler{
		jobQueue:      make(chan *Job, 100),
		priorityQueue: make(JobPriorityQueue, 0),
		workerCount:   workerCount,
		stopChan:      make(chan struct{}),
	}
	s.paused.Store(false)
	heap.Init(&s.priorityQueue)

	s.wg.Add(1)
	go s.dispatcher()

	s.wg.Add(s.workerCount)
	for i := 0; i < s.workerCount; i++ {
		go s.worker(i + 1)
	}

	fmt.Printf("Scheduler started with %d workers.\n", s.workerCount)
	return s
}

// dispatcher moves jobs from the priority queue to the job channel.
func (s *Scheduler) dispatcher() {
	defer s.wg.Done()
	fmt.Println("Dispatcher started.")
	for {
		// Check if paused first
		if s.paused.Load() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		select {
		case <-s.stopChan:
			close(s.jobQueue)
			fmt.Println("Dispatcher stopping.")
			return
		default:
			s.mu.Lock()
			if len(s.priorityQueue) > 0 {
				job := heap.Pop(&s.priorityQueue).(*Job)
				s.mu.Unlock()
				s.jobQueue <- job
			} else {
				s.mu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

// worker executes jobs from the job channel.
func (s *Scheduler) worker(id int) {
	defer s.wg.Done()
	fmt.Printf("Worker %d started\n", id)
	for job := range s.jobQueue {
		fmt.Printf("Worker %d starting job with priority %d\n", id, job.Priority)
		job.Task()
		fmt.Printf("Worker %d finished job\n", id)
	}
	fmt.Printf("Worker %d shutting down\n", id)
}

func (s *Scheduler) Pause() {
	s.paused.Store(true)
	fmt.Println("Scheduler paused.")
}

func (s *Scheduler) Resume() {
	s.paused.Store(false)
	fmt.Println("Scheduler resumed.")
}

// AddJob adds a new job to the priority queue.
func (s *Scheduler) AddJob(task func(), priority int) {
	job := &Job{Task: task, Priority: priority}
	s.mu.Lock()
	heap.Push(&s.priorityQueue, job)
	s.mu.Unlock()
}

// Stop shuts down the scheduler gracefully.
func (s *Scheduler) Stop() {
	// Signal the dispatcher to stop
	close(s.stopChan)
	// Wait for the dispatcher and all workers to finish their current tasks and exit
	s.wg.Wait()
	fmt.Println("Scheduler stopped.")
}

func main() {
	scheduler := NewScheduler(2)

	for i := 1; i <= 5; i++ {
		jobID := i
		scheduler.AddJob(func() {
			fmt.Printf("Executing job %d...\n", jobID)
			time.Sleep(2 * time.Second)
		}, 5)
	}

	fmt.Println("\n--- Letting 2 jobs run ---")
	time.Sleep(1 * time.Second)

	fmt.Println("\n--- PAUSING SCHEDULER ---")
	scheduler.Pause()
	fmt.Println("Scheduler will not dispatch new jobs for 4 seconds.")
	time.Sleep(4 * time.Second)

	fmt.Println("\n--- RESUMING SCHEDULER ---")
	scheduler.Resume()

	time.Sleep(5 * time.Second)

	scheduler.Stop()
}
