package util

import (
	"io/ioutil"
	"os"
)

// DirectoryExists checks if the directory path exists
func DirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// DeleteDirectory delete directory passed in path parameter
func DeleteDirectory(path string) error {
	// delete directory contents
	var err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

// ListDirectories returns a list of all directories in the given path parameter
func ListDirectories(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dirs := []string{}
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs, nil
}

// AddDirectory adds a new directory for the given path
func AddDirectory(path string) error {
	if !DirectoryExists(path) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	} else {
		err := os.Chmod(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
