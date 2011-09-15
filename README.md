Description
---------------
This is a project derived from [uniqush](http://uniqush.org).

With a task queue, a program can specify a task and its execution time. Then the task will be run at specified time point.

Install
---------------
`goinstall github.com/monnand/gotaskqueue`

Example
--------------
    package main

    import (
        "github.com/monnand/gotaskqueue"
        "fmt"
        "time"
    )

    // Create a new implementation of gotaskqueue.Task,
    // which must implement the Run() and ExecTime() methods.
    type MyTask struct {
        id int
        stop chan bool
        execTime int64
    }

    // Just a hello world.
    func (t *MyTask) Run(currentTime int64) {
        fmt.Printf("I am #%d task\n", t.id)
        t.stop <- true
    }

    // Return the executing time, in terms of
    // number of nanoseconds since the Unix epoch,
    // January 1, 1970 00:00:00 UTC.
    func (t *MyTask) ExecTime() int64 {
        return t.execTime
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
            // i + seconds. (i starts from 0)
            t.execTime = time.Nanoseconds() + (int64(i) + 1) * 1E9
            ch <- t
        }

        for i := 0; i < nr_tasks; i++ {
            <-stop
        }
    }

