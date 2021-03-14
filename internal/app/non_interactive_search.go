package app

import (
	"errors"

	"github.com/haroldadmin/getignore/internal/git"
)

func (app *GetIgnore) nonInteractiveSearch() (git.GitIgnoreFile, error) {
	results := app.ignores.SearchIgnores(app.searchQuery, 1)
	if len(results) == 0 {
		return git.GitIgnoreFile{}, errors.New("no matches found")
	}

	return results[0], nil
}
