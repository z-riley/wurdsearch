package main

import (
	"fmt"
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

// push puts an item at the back of the frontier queue. Returns error if the item already exists
func (f *frontier) push(item any) error {
	if f.contains(item) {
		return fmt.Errorf("Frontier already contains: %v", item)
	} else {
		err := f.queue.Enqueue(item)
		if err != nil {
			return err
		}
		return nil
	}
}

// contains checks if an item exists in the frontier queue
func (f *frontier) contains(item any) bool {
	for _, a := range f.queue.Slice {
		if a == item {
			return true
		}
	}
	return false
}

// getAll returns every item in the frontier as strings
func (f *frontier) getAll() ([]string, error) {
	// Make copy first for thread safety
	copy := append([]any{}, f.queue.Slice...)

	s, err := toStringSlice(copy)
	if err != nil {
		return []string{}, nil
	}

	return s, nil
}

type row struct {
	val   string
	count int
}

// TODO: topWebSites gets the most common websites in the frontier, ordered by count
func (f *frontier) countTopOccurrances() ([]row, error) {
	return []row{}, nil
}

// toStringSlice returns a copy of an any slice as a string slice
func toStringSlice(a []any) ([]string, error) {
	var s []string
	for _, v := range a {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("Non-string element found: %v", v)
		}
		s = append(s, str)
	}
	return s, nil
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
