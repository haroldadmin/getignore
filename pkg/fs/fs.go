package fs

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/go-git/go-billy/v5"
	"github.com/haroldadmin/getignore/pkg/utils"
)

var defaultSkipDirs []string = []string{".git"}
var logger = log.WithField("name", "fs")

var (
	ErrStatFailed = errors.New("stat-failed")
	ErrNotDir     = errors.New("not-a-directory")
	ErrReadDir    = errors.New("failed-to-read-dir")
)

type DiscoveredFile struct {
	Name      string
	Path      string
	IsSymlink bool
}

func ReadDirRecursively(
	filesystem billy.Filesystem,
	customSkipDirs ...string,
) ([]DiscoveredFile, error) {
	allFiles := []DiscoveredFile{}
	queue := utils.NewQueue()
	skipDirs := utils.NewSet(customSkipDirs...).Add(defaultSkipDirs...)

	readDir := func(path string) error {
		dirName := filepath.Base(path)
		if skipDirs.Contains(dirName) {
			logger.Debugf("Skipped: %s", dirName)
			return nil
		}

		dirEntries, err := filesystem.ReadDir(path)
		if err != nil {
			message := "failed to read directory"
			logger.Errorf("%s (%q): %v", message, path, err)
			return ErrReadDir
		}

		for _, entry := range dirEntries {
			if entry.IsDir() {
				logger.Debugf("Dir: %s", entry.Name())
				queue.Add(filepath.Join(path, entry.Name()))
				continue
			}

			if entry.Mode()&os.ModeSymlink == os.ModeSymlink {
				logger.Debugf("Symlink: %s", entry.Name())
				discoveredFile := DiscoveredFile{
					Name:      entry.Name(),
					Path:      filepath.Join(path, entry.Name()),
					IsSymlink: true,
				}
				allFiles = append(allFiles, discoveredFile)
				continue
			}

			logger.Debugf("File: %s", entry.Name())
			discoveredFile := DiscoveredFile{
				Name:      entry.Name(),
				Path:      filepath.Join(path, entry.Name()),
				IsSymlink: false,
			}
			allFiles = append(allFiles, discoveredFile)
		}

		return nil
	}

	readDir("/")

	for queue.Length() != 0 {
		first, err := queue.RemoveFirst()
		if err != nil {
			break
		}
		readDir(first)
	}

	return allFiles, nil
}
