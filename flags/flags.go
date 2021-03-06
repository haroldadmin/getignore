package flags

import "flag"

var SearchQuery string
var NoUpdate bool

func init() {
	flag.StringVar(&SearchQuery, "search", "", "Supply a search query to find gitignore files")
	flag.BoolVar(&NoUpdate, "no-update", false, "Controls if getignore should skip fetching repo updates")
	flag.Parse()
}
