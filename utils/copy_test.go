package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestCopyFile checks that a file will be copied from source to destination path
func TestCopyFile(t *testing.T) {
	srcFilePath := "/tmp/test-copy-file1"
	destFilePath := "/tmp/test-copy-file2"

	expectedText := "test\ncopy\nfile\n"
	err := ioutil.WriteFile(srcFilePath, []byte(expectedText), 0644)
	if err != nil {
		t.Fatalf("Сannot create src file: %s", err)
	}
	defer func() {
		os.Remove(srcFilePath)
	}()

	err = CopyFile(srcFilePath, destFilePath)
	if err != nil {
		t.Fatalf("Cannot copy file %s --> %s: %s", srcFilePath, destFilePath, err)
	}
	defer func() {
		os.Remove(destFilePath)
	}()

	bytes, err := ioutil.ReadFile(destFilePath)
	if err != nil {
		t.Fatalf("Сannot read dest file: %s", err)
	}

	if expectedText != string(bytes) {
		t.Fatal("Expected file content does not equal real")
	}
}

// TestCopyDir checks that a directory will be copied from source to destination path
func TestCopyDir(t *testing.T) {
	srcDir := "/tmp/test-copy-dir1/"
	destDir := "/tmp/test-copy-dir2/"

	fileName := "1.txt"
	childDir := "/child"
	filePath := childDir + "/" + fileName

	err := os.MkdirAll(srcDir+childDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Сannot create src dir: %s", err)
	}

	expectedText := "test\ncopy\ndir\n"
	err = ioutil.WriteFile(srcDir+filePath, []byte(expectedText), 0644)
	if err != nil {
		t.Fatalf("Сannot create src file: %s", err)
	}

	err = CopyDir(srcDir, destDir)
	defer func() {
		// clean data after test
		os.RemoveAll(srcDir)
		os.RemoveAll(destDir)
	}()
	if err != nil {
		t.Fatalf("Cannot copy dir: %s", err)
	}

	bytes, err := ioutil.ReadFile(destDir + filePath)
	if err != nil {
		t.Fatalf("Сannot read dest file: %s", err)
	}

	if expectedText != string(bytes) {
		t.Fatal("Expected file content does not equal real")
	}
}
