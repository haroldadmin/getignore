package main

import (
	"log"
	"strings"

	"github.com/haroldadmin/getignore/flags"
	_ "github.com/haroldadmin/getignore/flags"
	"github.com/haroldadmin/getignore/git"
	"github.com/manifoldco/promptui"
)

func main() {
	log.SetFlags(0)
	repoDir, err := git.Clone(flags.NoUpdate)
	if err != nil {
		log.Fatal(err)
	}

	ignores := &git.Ignores{
		RepoDir: repoDir,
	}

	err = ignores.FindIgnores()
	if err != nil {
		log.Fatal(err)
	}

	if flags.SearchQuery != "" {
		results := search(ignores, flags.SearchQuery)
		if len(results) == 0 {
			log.Printf("No matches found")
		}
		output := strings.Join(processMatches(results), ", ")
		log.Print(output)
		return
	}

	searchPrompt := promptui.Prompt{Label: "Search"}
	result, err := searchPrompt.Run()
	if err != nil {
		log.Printf("An error occurred while reading the prompt: %v", err)
		return
	}

	results := search(ignores, result)
	if len(results) == 0 {
		log.Printf("No matches found")
		return
	}

	selectionPrompt := promptui.Select{
		Label: "Results",
		Items: processMatches(results),
	}

	_, selection, err := selectionPrompt.Run()
	if err != nil {
		log.Printf("Failed to read selection result: %v", err)
		return
	}

	log.Print(selection)
}

func search(ignores *git.Ignores, query string) []git.GitIgnoreFile {
	matches := ignores.SearchIgnores(query, 5)
	if len(matches) == 0 {
		return []git.GitIgnoreFile{}
	}

	return matches
}

func processMatches(matches []git.GitIgnoreFile) []string {
	names := make([]string, 0, len(matches))
	for _, match := range matches {
		names = append(names, match.Name)
	}
	return names
}
