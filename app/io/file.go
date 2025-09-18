package io

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// WriteJSONFile will write v into absFilePath.
// v needs to be JSON marshalable.
func WriteJSONFile(absFilePath string, v any) error {
	content, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(absFilePath, content, 0644)
	return err
}

// DoesImageExistsOnDisk searches for the image on disk.
// First it will check for the absolute path (given via absImageFilePath).
// Then, we run a second method: Probalistic search. Mainly, because
// sometimes we get a value like "../generated/images/macht-der-craft"
// without any extension. This happens, e.g., when we don't get an
// image from the API. The probalistic search still checks if we have the
// image on disk.
func DoesImageExistsOnDisk(absImageFilePath string, onlyAbsoluteCheck bool) (string, bool) {
	// Check for an absolute match (full path)
	_, imageExistErr := os.Stat(absImageFilePath)
	if imageExistErr == nil {
		return absImageFilePath, true
	}

	// No fuzzy search
	if onlyAbsoluteCheck {
		return "", false
	}

	imagePath := path.Dir(absImageFilePath)
	imageFile := path.Base(absImageFilePath)

	// Sometimes, we don't get an empty image back
	// This means absImageFilePath is not having a file extension
	// Still, there might be a case that we have the image already.
	// Hence we run a probalistic check here
	imageFiles, err := GetAllFilesFromDirectoryWithExtensions(imagePath, GetImageExtensions())
	if err != nil {
		return "", false
	}

	for _, f := range imageFiles {
		if strings.HasPrefix(f.Name(), imageFile+".") {
			log.Printf("Image found via the probalistic way: %s", f.Name())
			return f.Name(), true
		}
	}

	return "", false
}

// CopyFile copies the contents of the file named src to dst.
// If dst does not exist, it will be created.
// If dst exists, it will be overwritten.
func CopyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	// Copy contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	// Flush writes to stable storage
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
