package main

import "log"

func main() {

	c, err := newCrawler()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range []string{
		"https://google.com/",
		"https://reddit.com/",
		"https://example.com/",
	} {
		c.frontier.queue.Enqueue(url)
	}

	// go func() {
	if err = c.crawlForever(); err != nil {
		log.Fatal(err)
	}
	// }()

}

func UNUSED(x ...any) {}
