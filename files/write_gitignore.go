package files

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/haroldadmin/getignore/git"
)

func WriteGitignore(selection git.GitIgnoreFile, output string) error {
	src, err := os.Open(selection.Path)
	if err != nil {
		return fmt.Errorf("Failed to read %v: %v", selection.Path, err)
	}
	defer src.Close()

	dest, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Failed to open/create %v: %v", output, err)
	}
	defer dest.Close()

	bytesCopied, err := io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("Failed to copy gitignore file: %v", err)
	}

	log.Printf("Copied %d bytes from %v to: %v", bytesCopied, selection.Path, output)
	return nil
}
