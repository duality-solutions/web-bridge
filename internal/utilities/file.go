package util

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// FileExists determines if the file exists for the given path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// MoveFile moves a file from source to destination
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

// DeleteFile delete the file passed in path parameter
func DeleteFile(path string) error {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// GetFileHash returns the hash for the given path parameter
func GetFileHash(path string) (string, error) {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return "", errOpen
	}
	defer f.Close()

	h := sha256.New()
	if _, errHash := io.Copy(h, f); errHash != nil {
		return "", errHash
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
