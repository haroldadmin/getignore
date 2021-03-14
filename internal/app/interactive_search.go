package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/haroldadmin/getignore/internal/git"
	"github.com/manifoldco/promptui"
)

func (app *GetIgnore) interactiveSearch(ctx context.Context) (git.GitIgnoreFile, error) {
	var file git.GitIgnoreFile

	// Loop while context is not cancelled or until we find a search result
	for ctx.Err() == nil {
		searchPrompt := promptui.Prompt{Label: "Search"}
		query, err := searchPrompt.Run()
		if err != nil {
			return file, fmt.Errorf("prompt error: %v", err)
		}

		results := app.ignores.SearchIgnores(query, 5)
		if len(results) == 0 {
			return file, errors.New("no matches found")
		}

		options := make([]string, 0, len(results))
		for _, result := range results {
			options = append(options, result.Name)
		}
		options = append(options, "search again")

		selectionPrompt := promptui.Select{
			Label: "Results",
			Items: options,
		}

		index, _, err := selectionPrompt.Run()
		if err != nil {
			return file, fmt.Errorf("failed to read selection result: %v", err)
		}

		if index == len(results) {
			// User wants to retry searching
			continue
		}

		return results[index], nil
	}

	return file, fmt.Errorf("context cancelled: %v", ctx.Err())
}
