package fsutil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ZipCreate(zipPath string, filePaths []string) (string, error) {
	if FileExists(zipPath) {
		zipPath = FileIncrementName(zipPath)
	}
	zipF, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipF.Close()

	zipW := zip.NewWriter(zipF)
	defer zipW.Close()

	for _, filePath := range filePaths {
		err = zipAddFile(filePath, zipW)
		if err != nil {
			return zipPath, err
		}
	}

	return zipPath, nil
}

func zipAddFile(filePath string, zipW *zip.Writer) error {
	originalF, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalF.Close()

	fileBase := filepath.Base(filePath)
	zipFileW, err := zipW.Create(fileBase)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFileW, originalF)
	if err != nil {
		return err
	}

	return nil
}

func ZipReadFileAsString(zipR *zip.ReadCloser, filePath string) (string, error) {
	var result string
	var zippedFile *zip.File

	for _, f := range zipR.File {
		if f.Name != filePath {
			continue
		}
		zippedFile = f
	}

	zippedFileR, err := zippedFile.Open()
	if err != nil {
		return result, err
	}
	defer zippedFileR.Close()

	bytes, err := io.ReadAll(zippedFileR)
	if err != nil {
		return result, err
	}

	result = string(bytes)

	return result, nil
}
