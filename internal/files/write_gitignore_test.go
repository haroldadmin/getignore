package files_test

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/haroldadmin/getignore/internal/files"
	"github.com/haroldadmin/getignore/internal/git"
)

func TestWriteGitignore(t *testing.T) {
	t.Run("should not produce errors when reading/writing from valid sources", func(t *testing.T) {
		randFilePath, inputCleanup, err := createInputFile(t)
		if err != nil {
			t.Fatal(err)
		}

		outputFilePath, outputCleanup, err := createOutputFile(t)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(inputCleanup)
		t.Cleanup(outputCleanup)

		gitIgnoreFile := git.GitIgnoreFile{
			Name: "Go.gitignore",
			Path: randFilePath,
		}

		err = files.WriteGitignore(gitIgnoreFile, outputFilePath)
		if err != nil {
			t.Errorf("expected no errors, got: %v", err)
			return
		}
	})

	t.Run("should copy contents of source to dest exactly", func(t *testing.T) {
		randFilePath, inputCleanup, err := createInputFile(t)
		if err != nil {
			t.Fatal(err)
		}

		outputFilePath, outputCleanup, err := createOutputFile(t)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(inputCleanup)
		t.Cleanup(outputCleanup)

		gitIgnoreFile := git.GitIgnoreFile{
			Name: "Go.gitignore",
			Path: randFilePath,
		}

		files.WriteGitignore(gitIgnoreFile, outputFilePath)

		inputContents, err := readFile(randFilePath, t)
		if err != nil {
			t.Fatal(err)
		}
		outputContents, err := readFile(outputFilePath, t)
		if err != nil {
			t.Fatal(err)
		}

		isContentEqual := reflect.DeepEqual(inputContents, outputContents)
		if !isContentEqual {
			t.Error("expected file contents to be the same, but were different")
			return
		}
	})
}

func createInputFile(t *testing.T) (string, func(), error) {
	t.Helper()

	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "gitignore-test-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp input file: %v", err)
	}

	defer tempFile.Close()

	writer := bufio.NewWriter(tempFile)
	for i := 0; i < 1024; i++ {
		randCharacterCode := 65 + rand.Intn(90-65)
		randByte := byte(randCharacterCode)
		writer.WriteByte(randByte)
	}

	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}

func createOutputFile(t *testing.T) (string, func(), error) {
	t.Helper()

	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "gitignore-test-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp output file: %v", err)
	}
	defer tempFile.Close()

	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}

func readFile(path string, t *testing.T) ([]byte, error) {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v for reading: %v", path, err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read from file %v: %v", path, err)
	}

	return bytes, nil
}
