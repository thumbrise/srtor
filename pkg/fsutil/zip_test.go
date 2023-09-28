package fsutil

import (
	"archive/zip"
	"os"
	"reflect"
	"testing"
)

func TestZipCreate(t *testing.T) {
	// Arrange
	const zipPath = "test.zip"
	const filePath = "test1.txt"
	const fileContent = "This is a test1 file"

	err := FileWrite(fileContent, filePath)
	if err != nil {
		t.Fatal(err)
		return
	}

	filePaths := []string{filePath}

	//Act
	_, err = ZipCreate(zipPath, filePaths)
	if err != nil {
		t.Fatal(err)
		return
	}

	//Assert
	zipR, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
		return
	}

	filePathsInZip := make([]string, 0)
	for _, file := range zipR.File {
		filePathsInZip = append(filePathsInZip, file.Name)
	}

	if !reflect.DeepEqual(filePaths, filePathsInZip) {
		t.Errorf("not all files copied in zip, original files { %v } zipped files { %v }", filePaths, filePathsInZip)
		return
	}

	for _, filePath := range filePaths {
		contentZipped, err := ZipReadFileAsString(zipR, filePath)
		if err != nil {
			t.Fatal(err)
			return
		}

		contentOgirinal, err := FileReadAsString(filePath)
		if err != nil {
			t.Fatal(err)
			return
		}

		if contentZipped != contentOgirinal {
			t.Errorf("Content zipped { %v } not equal content original { %v }", contentZipped, contentOgirinal)
			return
		}
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
		return
	}

	zipR.Close()

	err = os.Remove(zipPath)
	if err != nil {
		t.Fatal(err)
		return
	}
}
