Description
---------------
This project is a sub-project under [uniqush](http://uniqush.org).

With gotaskqueue, a program could define several tasks and process them separately at specific time points.

Install
---------------
`goinstall github.com/monnand/gotaskqueue`

Example
--------------

This example creates 5 tasks. Every second, one task will be processed.

    package main

    import (
        "github.com/monnand/gotaskqueue"
        "fmt"
    )

    // Define a new implementation of gotaskqueue.Task,
    // which must implement the Run() and ExecTime() methods.
    type MyTask struct {
        id int
        stop chan bool
        // TaskTime defines time-related operations,
        // so that you do not need to define your own ExecTime().
        // By composing TaskTime, you can use After(),
        // AfterNanoseconds() to specify the time point at which
        // the task will be processed by calling its Run()
        // method.
        gotaskqueue.TaskTime
    }

    // Just a hello world.
    func (t *MyTask) Run(currentTime int64) {
        // Print the task id and exit
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

            // The task with id i, will be processed after
            // i + 1 seconds. (i starts from 0).
            // Because we use gotaskqueue.TaskTime, it is easy to 
            // to use the After() method.
            t.After(int64(i) + 1)
            ch <- t
        }

        // wait all tasks to be processed.
        for i := 0; i < nr_tasks; i++ {
            <-stop
        }
    }

