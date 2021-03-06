package app

import (
	"context"
	"fmt"
	"log"

	"github.com/haroldadmin/getignore/files"
	"github.com/haroldadmin/getignore/git"
)

// GetIgnoreOptions contains options to configure GetIgnore
type GetIgnoreOptions struct {
	ShouldUpdate bool
	SearchQuery  string
	Output       string
}

// GetIgnore is the application struct
type GetIgnore struct {
	searchQuery string
	output      string
	ignores     git.Ignores
}

// Create creates an instance of GetIgnore using the given options
func Create(ctx context.Context, options GetIgnoreOptions) (*GetIgnore, error) {
	// Fetch the gitignore repository from Github
	repoDir, err := git.Clone(ctx, options.ShouldUpdate)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize getignore: %v", err)
	}

	// Find all the gitignore files in the repository
	ignores := git.Ignores{
		RepoDir: repoDir,
	}
	err = ignores.FindIgnores()
	if err != nil {
		return nil, fmt.Errorf("Failed to search for gitignore files: %v", err)
	}

	return &GetIgnore{
		ignores:     ignores,
		searchQuery: options.SearchQuery,
		output:      options.Output,
	}, nil
}

// Start initiates the search process for GetIgnore
func (app *GetIgnore) Start(ctx context.Context) {
	var selection git.GitIgnoreFile
	var err error

	if app.searchQuery != "" {
		selection, err = app.nonInteractiveSearch()
	} else {
		selection, err = app.interactiveSearch(ctx)
	}

	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("Selected: %v", selection.Name)
	if err := files.WriteGitignore(selection, app.output); err != nil {
		log.Print(err)
	}
}
