# Go Concurrent Job Scheduler

This project is a lightweight, concurrent job scheduler written in Go. It is designed to execute a set of tasks (represented by functions) concurrently while providing fine-grained control over execution, including job prioritization and scheduler state management.

## Features

-   **Fixed-Size Worker Pool**: Limits the number of concurrently executing jobs using a pool of goroutine workers, preventing resource exhaustion.
-   **Priority-Based Job Execution**: Ensures that jobs with a higher priority (represented by a lower integer value) are executed before lower-priority jobs.
-   **Scheduler Lifecycle Control**: Provides simple methods to `Pause`, `Resume`, and `Stop` the scheduler gracefully.

## Architecture Overview

The scheduler is built around a few core components that work together to manage and execute jobs efficiently.

1.  **Job**: The basic unit of work, containing a `Task` function and its `Priority`.
2.  **Priority Queue**: Instead of a simple FIFO channel, we use a min-heap (`container/heap`) to store incoming jobs. This ensures that when we need the next job, we can always retrieve the one with the highest priority in O(log n) time. This queue is protected by a `sync.Mutex` for concurrent access.
3.  **Dispatcher**: A central goroutine that acts as the "brain" of the scheduler. Its only job is to move the highest-priority item from the `Priority Queue` to the `Worker Channel`. It is also responsible for checking the scheduler's state (e.g., paused, stopped).
4.  **Worker Channel**: A standard buffered Go channel that the `Dispatcher` pushes jobs into.
5.  **Workers**: A fixed number of goroutines that continuously listen on the `Worker Channel`. When a job appears, a free worker picks it up and executes its task.

The data flow is as follows:

```
                                      +------------------+
                                      |   Worker 1       |
                                      +------------------+
                                             ^
                                             |
+--------+     +-----------------+     +---------------+     +------------------+     +------------------+
| User   | --> | scheduler.AddJob| --> | Priority Queue| --> |   Dispatcher  | --> |  Worker Channel  |
+--------+     +-----------------+     +---------------+     +---------------+     +------------------+
                                       (Heap, Mutex)                              (Buffered Chan)    |
                                                                                                     v
                                                                                               +------------------+
                                                                                               |   Worker 2       |
                                                                                               +------------------+
                                                                                                     ^
                                                                                                     | ... and so on
```

## Getting Started

### Prerequisites

-   Go 1.18 or later.

### Running the Example
1. Pick any file from 3 and rename it to `scheduler.go`.
2.  Run the application from your terminal:

```bash
go run scheduler.go
```

You will see output demonstrating workers picking up jobs according to priority, the scheduler pausing, and then resuming to finish the work before shutting down gracefully.

## API and Usage

Here is a simple example demonstrating how to use the scheduler.

```go
package main

import (
    "fmt"
    "time"
)

// Assume the scheduler code from the tutorial is in this package

func main() {
    // 1. Create a new scheduler with 2 concurrent workers.
    scheduler := NewScheduler(2)

    // 2. Add jobs with different priorities. Lower number = higher priority.
    // This high-priority job should run first.
    scheduler.AddJob(func() {
        fmt.Println("Executing highest priority job (P1)...")
        time.Sleep(2 * time.Second)
    }, 1)

    // Add some lower-priority jobs.
    for i := 1; i <= 3; i++ {
        jobID := i
        scheduler.AddJob(func() {
            fmt.Printf("Executing normal priority job (P5) #%d...\n", jobID)
            time.Sleep(2 * time.Second)
        }, 5)
    }

    fmt.Println("--> Jobs added. The two highest priority jobs will start.")
    time.Sleep(2500 * time.Millisecond) // Wait for the first two jobs to finish.

    // 3. Pause the scheduler. No new jobs will be dispatched.
    fmt.Println("\n--> PAUSING SCHEDULER...")
    scheduler.Pause()
    fmt.Println("Scheduler is paused. No new jobs will start for 3 seconds.")
    time.Sleep(3 * time.Second)

    // 4. Resume the scheduler.
    fmt.Println("\n--> RESUMING SCHEDULER...")
    scheduler.Resume()

    // Allow time for the remaining jobs to be processed.
    time.Sleep(5 * time.Second)

    // 5. Stop the scheduler gracefully.
    fmt.Println("\n--> STOPPING SCHEDULER...")
    scheduler.Stop()
    fmt.Println("Application finished.")
}

```

### Public API

-   `NewScheduler(workerCount int) *Scheduler`: Creates, initializes, and starts a new scheduler with the specified number of workers.
-   `(s *Scheduler) AddJob(task func(), priority int)`: Adds a new job to the priority queue.
-   `(s *Scheduler) Pause()`: Pauses the dispatcher, preventing new jobs from being sent to workers. Running jobs will complete.
-   `(s *Scheduler) Resume()`: Resumes a paused dispatcher.
-   `(s *Scheduler) Stop()`: Initiates a graceful shutdown of the dispatcher and all workers. It waits for all goroutines to finish.

## Future Improvements

-   **Panic Recovery**: Implement `recover()` within each worker to prevent a panicking job from crashing a worker goroutine.
-   **Context Propagation**: Allow passing a `context.Context` to jobs for handling cancellations and timeouts.
-   **Job Results and Errors**: Add a mechanism for jobs to return results or errors to the caller.