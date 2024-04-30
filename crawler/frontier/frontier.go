package frontier

import (
	"fmt"
	"net/url"
	"sort"

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
	if f.Contains(item) {
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

type Row struct {
	Val   string
	Count int
}

// TopNWebsites gets the most common websites in the frontier, ordered by number of occurances
func (f *Frontier) TopNWebsites(n int) ([]Row, error) {
	// Map to count the occurrence of each website
	counts := make(map[string]int)
	for _, item := range f.queue.Slice {
		link, ok := item.(string)
		if !ok {
			return []Row{}, fmt.Errorf("found non-string item in frontier: %v", item)
		}
		parsedUrl, err := url.Parse(link)
		if err != nil {
			return []Row{}, err
		}
		url := parsedUrl.Scheme + "://" + parsedUrl.Host
		counts[url]++
	}

	var rows []Row
	for website, count := range counts {
		rows = append(rows, Row{Val: website, Count: count})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Count > rows[j].Count
	})

	if n > len(rows) {
		n = len(rows)
	}
	return rows[:n], nil
}

// Contains checks if an item exists in the frontier queue
func (f *Frontier) Contains(item any) bool {
	for _, a := range f.queue.Slice {
		if a == item {
			return true
		}
	}
	return false
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
