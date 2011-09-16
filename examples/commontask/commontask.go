package main

import (
    tq "github.com/monnand/gotaskqueue"
    "fmt"
)

// We define a struct to do the job.
type MyJob struct {
    id int
    ch chan bool
}

// Everytime the job is run, it only print its id.
// The parameter t is used to control the time of next running.
// If we only want the job run once, then just ignore the
// parameter t.
// More details about using this parameter, see:
//  examples/expbackoff/expbackoff.go
func (j *MyJob) Run(t tq.TimeSetter) {
    fmt.Printf("Hello, I am #%d task.\n", j.id)
    j.ch <- true
}

func main() {
    stop := make(chan bool)

    // Create a channel to put task
    ch := make(chan tq.Task)

    // Create a new task queue
    q := tq.NewTaskQueue(ch)

    // Run this task queue in a separate goroutine
    go q.Run()

    // We create 5 tasks.
    // Every second, there will be one running task.
    for i := 0; i < 5; i++ {
        j := new(MyJob)
        j.id = i
        j.ch = stop

        // We need two parameters: the Runnable interface,
        // and the channel communicates with task queue.
        t := tq.NewCommonTask(j, ch)

        // Set the running time.
        // The task with id i will be run after i + 1 seconds.
        t.After(int64(i) + 1)
        ch <- t
    }
    for i := 0; i < 5; i++ {
        <-stop
    }
}


