package flags

import (
	"flag"
	"path/filepath"
)

var SearchQuery string
var NoUpdate bool
var OutputFile string

func init() {
	flag.StringVar(&SearchQuery, "search", "", "Supply a search query to find gitignore files")
	flag.BoolVar(&NoUpdate, "no-update", false, "Controls if getignore should skip fetching repo updates")
	flag.StringVar(&OutputFile, "out", filepath.Join(".", ".gitignore"), "The path of the .gitignore file to write")
	flag.Parse()
}
