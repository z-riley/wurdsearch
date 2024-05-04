package lemmatiser

import (
	"testing"
)

func TestLemmatise(t *testing.T) {
	l, err := NewLemmatiser()
	if err != nil {
		t.Fatal(err)
	}

	word := "cyclists"
	lemma := l.Lemmatise(word)
	if lemma != "cyclist" {
		t.Fail()
	}

	word = "anodised"
	lemma = l.Lemmatise(word)
	if lemma != "anodise" {
		t.Fail()
	}
}
