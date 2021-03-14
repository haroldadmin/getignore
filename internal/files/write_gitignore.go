package files

import (
	"fmt"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/haroldadmin/getignore/internal/git"
)

func WriteGitignore(selection git.GitIgnoreFile, output string) error {
	src, err := os.Open(selection.Path)
	if err != nil {
		return fmt.Errorf("failed to read %v: %v", selection.Path, err)
	}
	defer src.Close()

	dest, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to open/create %v: %v", output, err)
	}
	defer dest.Close()

	bytesCopied, err := io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("failed to copy gitignore file: %v", err)
	}

	log.Debugf("Copied %d bytes from %v to: %v", bytesCopied, selection.Path, output)
	return nil
}