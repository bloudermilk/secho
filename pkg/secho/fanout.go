package secho

import (
  "sync"
)

type Fanout struct {
  sync.Mutex
  Input chan Reading
  subscriptions []*chan Reading
}

func (fanout *Fanout) Subscribe () *chan Reading {
  channel := make(chan Reading)

  fanout.Lock()
  fanout.subscriptions = append(fanout.subscriptions, &channel)
  fanout.Unlock()

  return &channel
}

func (fanout *Fanout) DoFan () {
  for {
    reading := <-fanout.Input

    for _, subscription := range(fanout.subscriptions) {
      *subscription <- reading
    }
  }
}
