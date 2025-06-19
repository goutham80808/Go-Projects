// package main

// import (
// 	"container/heap"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Job struct {
// 	Task     func()
// 	Priority int
// }
// type JobPriorityQueue []*Job

// func (pq JobPriorityQueue) Len() int           { return len(pq) }
// func (pq JobPriorityQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }

// func (pq JobPriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

// func (pq *JobPriorityQueue) Push(x interface{}) {
// 	*pq = append(*pq, x.(*Job))
// }

// func (pq *JobPriorityQueue) Pop() interface{} {
// 	old := *pq
// 	n := len(old)
// 	item := old[n-1]
// 	(*pq)[n-1] = nil
// 	*pq = old[:n-1]
// 	return item
// }

// // Scheduler manages prioritized jobs and a worker pool.
// type Scheduler struct {
// 	jobQueue      chan *Job
// 	priorityQueue JobPriorityQueue
// 	mu            sync.Mutex
// 	workerCount   int
// 	wg            sync.WaitGroup
// 	stopChan      chan struct{}
// }

// // NewScheduler creates and starts the scheduler.
// func NewScheduler(workerCount int) *Scheduler {
// 	s := &Scheduler{
// 		jobQueue:      make(chan *Job, 100),
// 		priorityQueue: make(JobPriorityQueue, 0),
// 		workerCount:   workerCount,
// 		stopChan:      make(chan struct{}),
// 	}

// 	heap.Init(&s.priorityQueue)

// 	s.wg.Add(1)
// 	go s.dispatcher()

// 	s.wg.Add(s.workerCount)
// 	for i := 0; i < s.workerCount; i++ {
// 		go s.worker(i + 1)
// 	}

// 	fmt.Printf("Scheduler started with %d workers.\n", s.workerCount)
// 	return s
// }

// // dispatcher moves jobs from the priority queue to the job channel.
// func (s *Scheduler) dispatcher() {
// 	defer s.wg.Done()
// 	fmt.Println("Dispatcher started.")

// 	for {
// 		select {
// 		case <-s.stopChan:
// 			close(s.jobQueue)
// 			fmt.Println("Dispatcher stopping.")
// 			return
// 		default:
// 			s.mu.Lock()
// 			if len(s.priorityQueue) > 0 {
// 				job := heap.Pop(&s.priorityQueue).(*Job)
// 				s.mu.Unlock()
// 				s.jobQueue <- job
// 			} else {
// 				s.mu.Unlock()
// 				time.Sleep(10 * time.Millisecond)
// 			}
// 		}
// 	}
// }

// // worker executes jobs from the job channel.
// func (s *Scheduler) worker(id int) {
// 	defer s.wg.Done()
// 	fmt.Printf("Worker %d started\n", id)
// 	for job := range s.jobQueue {
// 		fmt.Printf("Worker %d starting job with priority %d\n", id, job.Priority)
// 		job.Task()
// 		fmt.Printf("Worker %d finished job\n", id)
// 	}
// 	fmt.Printf("Worker %d shutting down\n", id)
// }

// // AddJob adds a new job to the priority queue.
// func (s *Scheduler) AddJob(task func(), priority int) {
// 	job := &Job{Task: task, Priority: priority}
// 	s.mu.Lock()
// 	heap.Push(&s.priorityQueue, job)
// 	s.mu.Unlock()
// }

// // Stop shuts down the scheduler gracefully.
// func (s *Scheduler) Stop() {
// 	// Signal the dispatcher to stop
// 	close(s.stopChan)
// 	// Wait for the dispatcher and all workers to finish their current tasks and exit
// 	s.wg.Wait()
// 	fmt.Println("Scheduler stopped.")
// }

// func main() {
// 	scheduler := NewScheduler(3) // 3 workers

// 	tasks := map[int]int{
// 		1: 1,
// 		2: 5,
// 		3: 2,
// 		4: 10,
// 		5: 1,
// 		6: 7,
// 		7: 3,
// 	}

// 	for id, priority := range tasks {
// 		jobID := id
// 		jobPriority := priority

// 		scheduler.AddJob(func() {
// 			fmt.Printf("Executing job %d (Priority %d)...\n", jobID, jobPriority)
// 			time.Sleep(time.Duration(jobPriority) * 50 * time.Millisecond) // Simulate work based on priority
// 		}, jobPriority)
// 	}

// 	// Give some time for jobs to be processed
// 	time.Sleep(5 * time.Second)

// 	fmt.Println("Main: Signalling scheduler to stop...")
// 	scheduler.Stop()
// 	fmt.Println("Main: Scheduler stopped. Exiting.")
// }
