package main

import (
    tq "github.com/monnand/gotaskqueue"
    "fmt"
)

type MyJob struct {
    nr_exec_times int
    ch chan bool
}

func (j *MyJob) Run(t tq.TimeSetter) {
    if j.nr_exec_times >= 3 {
        fmt.Printf("I'm tired, I will stop.\n")
        t.Stop()
        j.ch <- true
        return
    }
    fmt.Printf("Hello, I have been executed %d time(s).\n", j.nr_exec_times)
    j.nr_exec_times++
}

func main() {
    stop := make(chan bool)

    // Create a channel to put task
    ch := make(chan tq.Task)

    // Create a new task queue
    q := tq.NewTaskQueue(ch)

    // Run this task queue in a separate goroutine
    go q.Run()
    j := new(MyJob)
    j.ch = stop

    t := tq.NewExpBackoffTask(j, ch, 2)

    // The task will be executed after 2 seconds
    t.After(3)

    ch <- t
    <-stop
}


