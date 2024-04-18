package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)


func TestGetAll(t *testing.T) {
	f := newFrontier()
	for _, url := range sampleCrawledUrls() {
		f.queue.Enqueue(url)
	}

	contents, err := f.getAll()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(contents, sampleCrawledUrls()) {
		t.Error("Expected out different to known input")
	}
}

func TestCountOccurrances(t *testing.T) {
	actual, err := countOccurrances(sampleCrawledUrls())
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(actual, map[string]int{
		"reddit.com":          62,
		"www.telegraph.co.uk": 2,
	}) {
		t.Error("Actaul not not equal expected")
	}
}

func sampleCrawledUrls() []string {
	return []string{
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
		"https://reddit.com/r/unitedkingdom/comments/1c69zby/jk_rowling_gets_apology_from_journalist_after/",
		"https://reddit.com/r/unitedkingdom/",
		"https://reddit.com/r/unitedkingdom/",
		"https://reddit.com/r/unitedkingdom/comments/1c69zby/jk_rowling_gets_apology_from_journalist_after/",
		"https://www.telegraph.co.uk/news/2024/04/16/jk-rowling-holocaust-denier-allegation-rivkah-brown-novara/",
		"https://www.telegraph.co.uk/news/2024/04/16/jk-rowling-holocaust-denier-allegation-rivkah-brown-novara/",
		"https://reddit.com/r/AskUK/comments/1c66aak/whats_the_worst_work_gift_youve_ever_received/",
		"https://reddit.com/r/AskUK/",
		"https://reddit.com/r/AskUK/",
		"https://reddit.com/r/AskUK/comments/1c66aak/whats_the_worst_work_gift_youve_ever_received/",
		"https://reddit.com/r/AskUK/comments/1c66aak/whats_the_worst_work_gift_youve_ever_received/",
		"https://reddit.com/r/london/comments/1c6c87e/please_help_find_my_missing_friend/",
		"https://reddit.com/r/london/",
		"https://reddit.com/r/london/",
		"https://reddit.com/r/london/comments/1c6c87e/please_help_find_my_missing_friend/",
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
