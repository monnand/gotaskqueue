/*
 * Copyright 2011 Nan Deng
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// With gotaskqueue, a program could define several tasks and execute them separately at specific time points.
package gotaskqueue

import (
    "github.com/petar/GoLLRB/llrb"
    "time"
)

// A task interface.
type Task interface {
    // This method will be called at specific time.
    // The parameter is current time when this method is executed,
    // in nanoseconds since Unix epoch.
    Run(time int64)

    // Return the executing time, in terms of
    // number of nanoseconds since the Unix epoch,
    // January 1, 1970 00:00:00 UTC.
    ExecTime() int64
}

// Usually, users do not want to define ExecTime() when they define a task.
// They may only want to tell the task queue: execute this task after 3 seconds.
// In this case, the user-defined task could composit this structure.
type TaskTime struct {
    execTime int64
}

func (t *TaskTime) ExecTime() int64 {
    return t.execTime
}

func (t *TaskTime) SetExecTime(nanosec int64) {
    t.execTime = nanosec
}

func (t *TaskTime) After(seconds int64) {
    t.execTime = (time.Seconds() + seconds) * 1E9
}

func (t *TaskTime) AfterNanoseconds(nanoseconds int64) {
    t.execTime = time.Nanoseconds() + nanoseconds
}

type TaskQueue struct {
    tree *llrb.Tree
    ch <-chan Task
    waitTime int64
}

const (
    MAX_TIME int64 = 0x0FFFFFFFFFFFFFFF
)

func taskBefore (a, b interface{}) bool{
    return a.(Task).ExecTime() < b.(Task).ExecTime()
}

func NewTaskQueue(ch <-chan Task) *TaskQueue {
    ret := new(TaskQueue)
    ret.tree = llrb.New(taskBefore)
    ret.ch = ch
    ret.waitTime = MAX_TIME
    return ret
}

func (t *TaskQueue) Run() {
    for {
        select {
        case task := <-t.ch:
            if task == nil {
                return
            }
            /*
            fmt.Printf("I received a task. Current Time %d, it ask me to run it at %d\n",
                        time.Nanoseconds()/1E9, task.ExecTime()/1E9)
                        */
            if task.ExecTime() <= time.Nanoseconds() {
                go task.Run(time.Nanoseconds())
                continue
            }
            t.tree.ReplaceOrInsert(task)
            x := t.tree.Min()
            if x == nil {
                t.waitTime = MAX_TIME
                continue
            }
            task = x.(Task)
            t.waitTime = task.ExecTime() - time.Nanoseconds()
        case <-time.After(t.waitTime):
            x := t.tree.Min()
            if x == nil {
                t.waitTime = MAX_TIME
                continue
            }
            task := x.(Task)
            /*
            fmt.Printf("Current Time %d, a task ask me to run it at %d\n",
                        time.Nanoseconds()/1E9, task.ExecTime()/1E9)
                        */
            for task.ExecTime() <= time.Nanoseconds() {
                go task.Run(time.Nanoseconds())
                t.tree.DeleteMin()
                x = t.tree.Min()
                if x == nil {
                    t.waitTime = MAX_TIME
                    task = nil
                    break
                }
                /*
                fmt.Printf("Current Time %d, a task ask me to run it at %d\n",
                        time.Nanoseconds()/1E9, task.ExecTime()/1E9)
                        */
                task = x.(Task)
            }
            if task != nil {
                t.waitTime = task.ExecTime() - time.Nanoseconds()
            }
        }
    }
}

