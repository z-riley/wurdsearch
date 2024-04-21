package main

// website contains data about which other sites link to and from itself.
// This is used to help calculate the "importance" of the site
type website struct {
	url       string   `bson:"url"`
	linksTo   []string `bson:"linksTo"`
	linksFrom []string `bson:"linksFrom"`
}
