package main

import (
	"github.com/enriquebris/goconcurrentqueue"
)

type frontier struct {
	queue *goconcurrentqueue.FIFO
}

func newFrontier() *frontier {
	return &frontier{
		queue: goconcurrentqueue.NewFIFO(),
	}
}
