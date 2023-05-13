package processor

import (
	hqurl "github.com/hueristiq/hqgoutils/url"
)

// a Extractor is any function that accepts a URL and some
// kind of format string (which may not actually be used
// by some functions), and returns a slice of strings
// derived from that URL. It's not uncommon for a Extractor
// function to return a slice of length 1, but the return
// type remains a slice because *some* functions need to
// return multiple strings; e.g. the keys function.
type Extractor func(*hqurl.URL, string) []string
