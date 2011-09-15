package main

import (
    "github.com/monnand/gotaskqueue"
    "fmt"
)

// Create a new implementation of gotaskqueue.Task,
// which must implement the Run() and ExecTime() methods.
type MyTask struct {
    id int
    stop chan bool
    // See details below
    gotaskqueue.TaskTime
}

// Just a hello world.
func (t *MyTask) Run(currentTime int64) {
    fmt.Printf("I am #%d task\n", t.id)
    t.stop <- true
}

func main() {
    stop := make(chan bool)
    nr_tasks := 5

    // Create a channel to put task
    ch := make(chan gotaskqueue.Task)

    // Create a new task queue
    q := gotaskqueue.NewTaskQueue(ch)

    // Run this task queue in a separate goroutine
    go q.Run()

    // Insert 5 tasks into the queue.
    for i := 0; i < nr_tasks; i++ {
        t := new(MyTask)
        t.id = i
        t.stop = stop

        // The execution time of the task.
        // The task with id i, will be executed after
        // i + seconds. (i starts from 0).
        // Because we use gotaskqueue.TaskTime, it is easy to 
        // to use the After() method.
        t.After(int64(i) + 1)
        ch <- t
    }

    for i := 0; i < nr_tasks; i++ {
        <-stop
    }
}

