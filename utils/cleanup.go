package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DeleteEmptyDirs(dir string, root string) error {

	dir = filepath.Clean(dir)
	root = filepath.Clean(root)

	for {
		if dir == root || len(dir) < len(root) {
			break
		}

		isEmpty, err := isDirEmpty(dir)
		if err != nil {
			return err
		}
		if !isEmpty {
			break
		}

		if err := os.Remove(dir); err != nil {
			return err
		}

		fmt.Println("Deleted empty dir: ", dir)

		dir = filepath.Dir(dir)
	}
	return nil
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
