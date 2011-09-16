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
    fmt.Printf("Hello, I have been executed %d time(s).\n", j.nr_exec_times)

    // After 5 times run, it will stop.
    if j.nr_exec_times >= 5 {
        t.Stop()
        j.ch <- true
        return
    }
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

    t := tq.NewPeriodicTask(j, ch)

    // The task will be executed after 2 seconds
    t.After(2)

    // After the first run, every second, this task will be executed again.
    t.SetPeriod(1)

    ch <- t
    <-stop
}

