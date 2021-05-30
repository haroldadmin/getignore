package gitignore

type GitIgnores []GitIgnoreFile

func (g GitIgnores) String(i int) string {
	return g[i].Name
}

func (g GitIgnores) Len() int {
	return len(g)
}
