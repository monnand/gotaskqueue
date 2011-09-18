package gotaskqueue

import (
    "time"
)

type TimeSetter interface {
    After(s int64)
    AfterNanoseconds(ns int64)
    Stop()
}

type Runnable interface {
    Run(s TimeSetter)
}

type CommonTask struct {
    task Runnable
    execTime int64
    ch chan<- Task
}

func (t *CommonTask) After(seconds int64) {
    t.execTime = (time.Seconds() + seconds) * 1E9
}

func (t *CommonTask) AfterNanoseconds(nanoseconds int64) {
    t.execTime = time.Nanoseconds() + nanoseconds
}

func (t *CommonTask) ExecTime() int64 {
    return t.execTime
}

func (t *CommonTask) Run(ct int64) {
    if t.task != nil {
        t.task.Run(t)
    }
    if t.execTime > ct {
        t.ch <- t
    }
}

func (t *CommonTask) Stop() {
    t.execTime = -1
    t.task = nil
}

func NewCommonTask(task Runnable, ch chan<- Task) *CommonTask {
    ret := new(CommonTask)
    ret.task = task
    ret.ch = ch
    return ret
}

type PeriodicTask struct {
    period int64
    CommonTask
}

func (t *PeriodicTask) Run(ct int64) {
    t.execTime = ct + t.period
    if t.task != nil {
        t.task.Run(t)
    } else {
        return
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

func NewPeriodicTask(task Runnable, ch chan<- Task) *PeriodicTask {
    ret := new(PeriodicTask)
    ret.task = task
    ret.period = -1
    ret.ch = ch
    return ret
}

type ExpBackoffTask struct {
    backoff int64
    CommonTask
}

func NewExpBackoffTask(task Runnable, ch chan<- Task, backoff int64) *ExpBackoffTask {
    ret := new(ExpBackoffTask)
    ret.task = task
    ret.backoff = backoff * 1E9
    ret.ch = ch
    return ret
}

func (t *ExpBackoffTask) Run(ct int64) {
    if t.task != nil {
        t.task.Run(t)
    } else {
        return
    }
    if t.backoff <= 0 {
        return
    }
    t.execTime = ct + t.backoff
    t.backoff <<= 1
    t.ch <- t
    return
}

