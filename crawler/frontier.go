package main

import (
	"github.com/enriquebris/goconcurrentqueue"
)

type Frontier struct {
	queue *goconcurrentqueue.FIFO
}

func newFrontier() *Frontier {
	return &Frontier{
		queue: goconcurrentqueue.NewFIFO(),
	}
}
