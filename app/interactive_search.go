package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/haroldadmin/getignore/git"
	"github.com/manifoldco/promptui"
)

func (app *GetIgnore) interactiveSearch(ctx context.Context) (git.GitIgnoreFile, error) {
	var file git.GitIgnoreFile

	// Loop while context is not cancelled or until we find a search result
	for ctx.Err() == nil {
		searchPrompt := promptui.Prompt{Label: "Search"}
		query, err := searchPrompt.Run()
		if err != nil {
			return file, fmt.Errorf("Prompt error: %v", err)
		}

		results := app.ignores.SearchIgnores(query, 5)
		if len(results) == 0 {
			return file, errors.New("No matches found")
		}

		options := make([]string, 0, len(results))
		for _, result := range results {
			options = append(options, result.Name)
		}
		options = append(options, "Search again")

		selectionPrompt := promptui.Select{
			Label: "Results",
			Items: options,
		}

		index, _, err := selectionPrompt.Run()
		if err != nil {
			return file, fmt.Errorf("Failed to read selection result: %v", err)
		}

		if index == len(results) {
			// User wants to retry searching
			continue
		}

		return results[index], nil
	}

	return file, fmt.Errorf("Context cancelled: %v", ctx.Err())
}
