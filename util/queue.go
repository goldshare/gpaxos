package util

import (
  "sync"
  "time"
)

type Queue struct {
  TimeoutCond *TimeoutCond
  Mutex       sync.Mutex
  Storage     *Deque
}

func NewQueue() *Queue {
  queue := &Queue{}

  queue.TimeoutCond = NewTimeoutCondWithMutex(&queue.Mutex)
  queue.Storage = NewDeque()

  return queue
}

func (self *Queue) Peek() interface{} {
  for self.Empty() {
    self.TimeoutCond.Wait()
  }

  return self.Storage.Pop()
}

func (self *Queue) PeekWithTimeout(timeout int) interface{} {
  for self.Empty() {
    self.TimeoutCond.Lock()
    ret := self.TimeoutCond.WaitFor(timeout)
    self.TimeoutCond.Unlock()
    if !ret {
      return nil
    }
  }

  return self.Storage.Pop()
}

func (self *Queue) Pop() {
  self.Storage.Pop()
}

func (self *Queue) Add(value interface{}) int {
  self.Storage.Append(value)

  self.TimeoutCond.Signal()

  return self.Storage.Size()
}

func (self *Queue) Size() int {
  return self.Storage.Size()
}

func (self *Queue) Signal() {
  self.TimeoutCond.Signal()
}

func (self *Queue) Broadcast() {
  self.TimeoutCond.Broadcast()
}

func (self *Queue) Lock() {
  self.Mutex.Lock()
}

func (self *Queue) Unlock() {
  self.Mutex.Unlock()
}

func (self *Queue) Empty() bool {
  return self.Storage.Empty()
}
