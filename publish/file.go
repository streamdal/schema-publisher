package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createZip(dir string) ([]byte, error) {
	newZipFile := &bytes.Buffer{}
	zipWriter := zip.NewWriter(newZipFile)

	if dir != "./" && dir != "." {
		previousDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("unable to get current dir: %s", err)
		}

		// Temporarily change working dir
		if err := os.Chdir(dir); err != nil {
			return nil, fmt.Errorf("unable to change dir to '%s': %s", dir, err)
		}

		defer os.Chdir(previousDir)
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %s", err)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("unable to open path '%s': %s", path, err)
		}

		defer file.Close()

		f, err := zipWriter.Create(path)
		if err != nil {
			return fmt.Errorf("unable to create zip entry for path '%s': %s", path, err)
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return fmt.Errorf("unable to copy file contents for '%s' to zip: %s", path, err)
		}

		return nil
	}

	if err := filepath.Walk(".", walker); err != nil {
		return nil, fmt.Errorf("unable to walk dir '%s': %s", dir, err)
	}

	zipWriter.Close()

	return newZipFile.Bytes(), nil
}
