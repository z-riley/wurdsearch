package stopwords

// stopWords is a slice of common stop words derived from various online sources
var stopWords = []string{
	"a",
	"about",
	"an",
	"and",
	"are",
	"as",
	"at",
	"be",
	"but",
	"by",
	"for",
	"from",
	"i",
	"if",
	"in",
	"into",
	"is",
	"it",
	"no",
	"not",
	"of",
	"on",
	"or",
	"such",
	"that",
	"the",
	"their",
	"then",
	"there",
	"these",
	"they",
	"this",
	"to",
	"was",
	"what",
	"when",
	"where",
	"who",
	"will",
	"with",
}

func IsStopWord(word string) bool {
	for _, stopWord := range stopWords {
		if word == stopWord {
			return true
		}
	}
	return false
}
