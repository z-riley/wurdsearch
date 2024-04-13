package main

import (
	"github.com/enriquebris/goconcurrentqueue"
)

type Frontier struct {
	Queue *goconcurrentqueue.FIFO
}

func newFrontier() *Frontier {
	return &Frontier{
		Queue: goconcurrentqueue.NewFIFO(),
	}
}
