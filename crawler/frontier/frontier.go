package frontier

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/enriquebris/goconcurrentqueue"
)

type Frontier struct {
	queue *goconcurrentqueue.FIFO
}

func NewFrontier() *Frontier {
	return &Frontier{
		queue: goconcurrentqueue.NewFIFO(),
	}
}

// Push puts an item at the back of the frontier queue. Returns error if the item already exists
func (f *Frontier) Push(item any) error {
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

// Dequeue removes an element from the front of the frontier. Returns error if queue is locked or empty
func (f *Frontier) Dequeue() (any, error) {
	return f.queue.Dequeue()
}

// GetAll returns every item in the frontier as strings
func (f *Frontier) GetAll() ([]string, error) {
	// Make copy first for thread safety
	copy := append([]any{}, f.queue.Slice...)

	s, err := toStringSlice(copy)
	if err != nil {
		return []string{}, nil
	}

	return s, nil
}

// GetLen returns the number elements in the frontier
func (f *Frontier) Len() int {
	return f.queue.GetLen()
}

// Contains returns true if frontier contains the a given item
func (f *Frontier) Contains(item any) bool {
	return slices.Contains(f.queue.Slice, item)
}

type Row struct {
	Val   string
	Count int
}

// TODO: topWebSites gets the most common websites in the frontier, ordered by count
func (f *Frontier) CountTopOccurrances() ([]Row, error) {
	return []Row{}, nil
}

// contains checks if an item exists in the frontier queue
func (f *Frontier) contains(item any) bool {
	for _, a := range f.queue.Slice {
		if a == item {
			return true
		}
	}
	return false
}

// CountOccurrances returns a map of the most common strings in a slice and their respective counts
func CountOccurrances(urls []string) (map[string]int, error) {
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
