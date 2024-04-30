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

	if err := f.Push("d"); err != nil {
		t.Fatal(err)
	}

	err := f.Push("a")
	if err == nil {
		t.Error("Should have errored when duplicate item pushed")
	}
}

func TestGetAll(t *testing.T) {
	f := NewFrontier()
	for _, url := range sampleCrawledUrls() {
		f.queue.Enqueue(url)
	}

	contents, err := f.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(contents, sampleCrawledUrls()) {
		t.Error("Expected out different to known input")
	}
}

func TestTopNWebsites(t *testing.T) {
	f := NewFrontier()
	for _, item := range sampleCrawledUrls() {
		if err := f.Push(item); err != nil {
			t.Fatal(err)
		}
	}

	result, err := f.TopNWebsites(10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

}

func sampleCrawledUrls() []string {
	return []string{
		"https://rue.wikipedia.org",
		"https://sah.wikipedia.org",
		"https://sat.wikipedia.org",
		"https://reddit.com/r/popular/hot/?geo_filter=us",
		"https://reddit.com/r/popular/hot/?geo_filter=ar",
		"https://reddit.com/r/popular/hot/?geo_filter=au",
		"https://reddit.com/r/popular/hot/?geo_filter=bg",
		"https://reddit.com/r/popular/hot/?geo_filter=ca",
		"https://reddit.com/r/popular/hot/?geo_filter=cl",
		"https://reddit.com/r/popular/hot/?geo_filter=co",
		"https://reddit.com/r/popular/hot/?geo_filter=hr",
		"https://reddit.com/r/popular/hot/?geo_filter=cz",
		"https://reddit.com/r/popular/hot/?geo_filter=fi",
		"https://reddit.com/r/popular/hot/?geo_filter=fr",
		"https://reddit.com/r/popular/hot/?geo_filter=de",
		"https://reddit.com/r/popular/hot/?geo_filter=gr",
		"https://reddit.com/r/popular/hot/?geo_filter=hu",
		"https://reddit.com/r/popular/hot/?geo_filter=is",
		"https://reddit.com/r/popular/hot/?geo_filter=in",
		"https://reddit.com/r/popular/hot/?geo_filter=ie",
		"https://reddit.com/r/popular/hot/?geo_filter=it",
		"https://reddit.com/r/popular/hot/?geo_filter=jp",
		"https://reddit.com/r/popular/hot/?geo_filter=my",
		"https://reddit.com/r/popular/hot/?geo_filter=mx",
		"https://reddit.com/r/popular/hot/?geo_filter=nz",
		"https://reddit.com/r/popular/hot/?geo_filter=ph",
		"https://reddit.com/r/popular/hot/?geo_filter=pl",
		"https://reddit.com/r/popular/hot/?geo_filter=pt",
		"https://reddit.com/r/popular/hot/?geo_filter=pr",
		"https://reddit.com/r/popular/hot/?geo_filter=ro",
		"https://reddit.com/r/popular/hot/?geo_filter=rs",
		"https://reddit.com/r/popular/hot/?geo_filter=sg",
		"https://reddit.com/r/popular/hot/?geo_filter=es",
		"https://reddit.com/r/popular/hot/?geo_filter=se",
		"https://reddit.com/r/popular/hot/?geo_filter=tw",
		"https://reddit.com/r/popular/hot/?geo_filter=th",
		"https://reddit.com/r/popular/hot/?geo_filter=tr",
		"https://reddit.com/r/popular/hot/?geo_filter=gb",
		"https://reddit.com/?feedViewType=cardView",
		"https://reddit.com/?feedViewType=compactView",
		"https://reddit.com/r/unitedkingdom/",
		"https://reddit.com/r/unitedkingdom/comments/1c69zby/jk_rowling_gets_apology_from_journalist_after/",
		"https://www.telegraph.co.uk/news/2024/04/16/jk-rowling-holocaust-denier-allegation-rivkah-brown-novara/",
		"https://reddit.com/r/AskUK/",
		"https://reddit.com/r/AskUK/comments/1c66aak/whats_the_worst_work_gift_youve_ever_received/",
		"https://reddit.com/r/london/comments/1c6c87e/please_help_find_my_missing_friend/",
		"https://reddit.com/r/london/",
		"https://reddit.com/r/DestinyTheGame",
		"https://reddit.com/r/anime",
		"https://reddit.com/r/destiny2",
		"https://reddit.com/r/FortNiteBR",
		"https://reddit.com/r/dndnext",
		"https://reddit.com/r/buildapc",
		"https://reddit.com/r/techsupport",
		"https://reddit.com/r/jailbreak",
		"https://reddit.com/r/LivestreamFail",
		"https://reddit.com/r/legaladvice",
		"https://reddit.com/r/LifeProTips",
		"https://reddit.com/r/AskCulinary",
	}
}
