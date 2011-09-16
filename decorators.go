package gotaskqueue

import (
    "time"
)

type Runnable interface {
    Run()
}

type OneTimeTask struct {
    task Runnable
    execTime int64
}

func (t *OneTimeTask) After(seconds int64) {
    t.execTime = (time.Seconds() + seconds) * 1E9
}

func (t *OneTimeTask) AfterNanoseconds(nanoseconds int64) {
    t.execTime = time.Nanoseconds() + nanoseconds
}

func (t *OneTimeTask) ExecTime() int64 {
    return t.execTime
}

func (t *OneTimeTask) Run(ct int64) {
    if t.task != nil {
        t.task.Run()
    }
}

func NewOneTimeTask(task Runnable) Task {
    ret := new(OneTimeTask)
    ret.task = task
    return ret
}

type PeriodicTask struct {
    task Runnable
    period int64
    execTime int64
    ch chan<- Task
}

func (t *PeriodicTask) ExecTime() int64 {
    return t.execTime
}

func (t *PeriodicTask) Run(ct int64) {
    t.execTime = ct + t.period
    if t.task != nil {
        t.task.Run()
    }
    if t.period <= 0 {
        return
    }
    t.ch <- t
    return
}

func (t *PeriodicTask) SetPeriod(seconds int64) {
    t.period = seconds * 1E9
}

func (t *PeriodicTask) SetPeriodInNanoseconds(ns int64) {
    t.period = ns
}

func NewPeriodicTask(task Runnable, ch chan<- Task) Task {
    ret := new(PeriodicTask)
    ret.task = task
    ret.period = -1
    ret.ch = ch
    return ret
}

