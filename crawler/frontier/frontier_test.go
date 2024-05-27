package frontier

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPush(t *testing.T) {
	f := NewFrontier()
	for _, url := range []string{
		"a",
		"b",
		"c",
	} {
		f.queue.Enqueue(url)
	}

	// Duplicate links
	if err := f.Push(Link{"d", 0}); err != nil {
		t.Fatal(err)
	}
	err := f.Push(Link{"d", 0})
	if err == nil {
		t.Error("Should have errored when duplicate item pushed")
	}
}

func TestGetAll(t *testing.T) {
	f := NewFrontier()
	for _, url := range sampleCrawledLinks() {
		f.queue.Enqueue(url)
	}

	contents := f.GetAll()

	if !reflect.DeepEqual(contents, sampleCrawledLinks()) {
		t.Error("Expected out different to known input")
	}
}

func TestTopNWebsites(t *testing.T) {
	f := NewFrontier()
	for _, link := range sampleCrawledLinks() {
		if err := f.Push(link); err != nil {
			t.Fatal(err)
		}
	}

	result, err := f.TopNWebsites(10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

}

func sampleCrawledLinks() []Link {
	return []Link{
		{"https://rue.wikipedia.org", 0},
		{"https://sah.wikipedia.org", 0},
		{"https://sat.wikipedia.org", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=us", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=ar", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=au", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=bg", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=ca", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=cl", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=co", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=hr", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=cz", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=fi", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=fr", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=de", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=gr", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=hu", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=is", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=in", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=ie", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=it", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=jp", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=my", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=mx", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=nz", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=ph", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=pl", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=pt", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=pr", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=ro", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=rs", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=sg", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=es", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=se", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=tw", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=th", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=tr", 0},
		{"https://reddit.com/r/popular/hot/?geo_filter=gb", 0},
		{"https://reddit.com/?feedViewType=cardView", 0},
		{"https://reddit.com/?feedViewType=compactView", 0},
		{"https://reddit.com/r/unitedkingdom/", 0},
		{"https://reddit.com/r/unitedkingdom/comments/1c69zby/jk_rowling_gets_apology_from_journalist_after/", 0},
		{"https://www.telegraph.co.uk/news/2024/04/16/jk-rowling-holocaust-denier-allegation-rivkah-brown-novara/", 0},
		{"https://reddit.com/r/AskUK/", 0},
		{"https://reddit.com/r/AskUK/comments/1c66aak/whats_the_worst_work_gift_youve_ever_received/", 0},
		{"https://reddit.com/r/london/comments/1c6c87e/please_help_find_my_missing_friend/", 0},
		{"https://reddit.com/r/london/", 0},
		{"https://reddit.com/r/DestinyTheGame", 0},
		{"https://reddit.com/r/anime", 0},
		{"https://reddit.com/r/destiny2", 0},
		{"https://reddit.com/r/FortNiteBR", 0},
		{"https://reddit.com/r/dndnext", 0},
		{"https://reddit.com/r/buildapc", 0},
		{"https://reddit.com/r/techsupport", 0},
		{"https://reddit.com/r/jailbreak", 0},
		{"https://reddit.com/r/LivestreamFail", 0},
		{"https://reddit.com/r/legaladvice", 0},
		{"https://reddit.com/r/LifeProTips", 0},
		{"https://reddit.com/r/AskCulinary", 0},
	}
}
