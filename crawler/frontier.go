package main

import (
	"net/url"

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

// TODO: make sure new links sent to frontier aren't already in the frontier

// getAll returns every item in the frontier as a string. Warning: this function isn't thread safe
func (f *frontier) getAll() ([]string, error) {
	var items []string

	for i := 0; i < f.queue.GetLen(); i++ {
		item, err := f.queue.Get(i)
		if err != nil {
			return items, err
		}

		items = append(items, item.(string))
	}

	return items, nil
}

// countOccurrances returns a map of the most common strings in a slice and their respective counts
func countOccurrances(urls []string) (map[string]int, error) {
	result := make(map[string]int)
	for _, URL := range urls {
		parsedUrl, err := url.Parse(URL)
		if err != nil {
			return result, err
		}
		result[parsedUrl.Host] += 1
	}
	return result, nil
}
