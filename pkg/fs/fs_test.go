package fs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/haroldadmin/getignore/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func TestReadDirRecursively(t *testing.T) {
	t.Run("it should read all top level files in the directory", func(t *testing.T) {
		dir := t.TempDir()
		fileCount := 10
		for i := 0; i < fileCount; i++ {
			testFile(t, dir, fmt.Sprint(i))
		}

		defer os.RemoveAll(dir)

		osFs := osfs.New(dir)
		files, err := fs.ReadDirRecursively(osFs)

		assert.NoError(t, err)
		assert.Equal(t, fileCount, len(files))
	})

	t.Run("it should read all nested files in the directory", func(t *testing.T) {
		parentDir := t.TempDir()
		nestedDirCount := 10
		nestedFileCount := 10
		for i := 0; i < nestedDirCount; i++ {
			nestedDir, cleanup := testDir(t, parentDir, fmt.Sprint(i))
			defer cleanup()

			for j := 0; j < nestedFileCount; j++ {
				testFile(t, nestedDir, fmt.Sprint(j))
			}
		}

		osFs := osfs.New(parentDir)
		files, err := fs.ReadDirRecursively(osFs)

		assert.NoError(t, err)
		assert.Equal(t, nestedDirCount*nestedFileCount, len(files))
	})

	t.Run("it should skip specified directories", func(t *testing.T) {
		parentDir := t.TempDir()
		nestedDirCount := 10
		nestedFileCount := 10
		skipSet := []string{}
		for i := 0; i < nestedDirCount; i++ {
			nestedDir, cleanup := testDir(t, parentDir, fmt.Sprint(i))
			defer cleanup()

			skipSet = append(skipSet, filepath.Base(nestedDir))

			for j := 0; j < nestedFileCount; j++ {
				testFile(t, nestedDir, fmt.Sprint(j))
			}
		}

		osFs := osfs.New(parentDir)
		files, err := fs.ReadDirRecursively(osFs, skipSet...)

		assert.NoError(t, err)
		assert.Empty(t, files)
	})
}

func testFile(t *testing.T, dir, fileName string) *os.File {
	t.Helper()
	path := filepath.Join(dir, fileName)
	file, _ := os.Create(path)
	return file
}

func testDir(t *testing.T, dir, dirName string) (string, func()) {
	t.Helper()
	path := filepath.Join(dir, dirName)
	os.MkdirAll(path, os.ModePerm)

	cleanup := func() { os.RemoveAll(path) }

	return path, cleanup
}
