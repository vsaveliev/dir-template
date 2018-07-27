package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestCopyFile checks that a file will be copied from source to destination path
func TestCopyFile(t *testing.T) {
	src := "/tmp/test-copy-file1"
	dest := "/tmp/test-copy-file2"

	expectedText := "test\ncopy\nfile\n"
	err := ioutil.WriteFile(src, []byte(expectedText), 0644)
	if err != nil {
		t.Fatalf("Сannot create src file: %s", err)
	}
	defer func() {
		os.Remove(src)
	}()

	err = CopyFile(src, dest)
	if err != nil {
		t.Fatalf("Cannot copy file %s --> %s: %s", src, dest, err)
	}
	defer func() {
		os.Remove(dest)
	}()

	bytes, err := ioutil.ReadFile(dest)
	if err != nil {
		t.Fatalf("Сannot read dest file: %s", err)
	}

	if expectedText != string(bytes) {
		t.Fatal("Expected file content does not equal real")
	}
}

// TestCopyDir checks that a directory will be copied from source to destination path
func TestCopyDir(t *testing.T) {
	src := "/tmp/test-copy-dir1/"
	dest := "/tmp/test-copy-dir2/"

	fileName := "1.txt"
	childDir := "/child"
	filePath := childDir + "/" + fileName

	err := os.MkdirAll(src+childDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Сannot create src dir: %s", err)
	}

	expectedText := "test\ncopy\ndir\n"
	err = ioutil.WriteFile(src+filePath, []byte(expectedText), 0644)
	if err != nil {
		t.Fatalf("Сannot create src file: %s", err)
	}

	err = CopyDir(src, dest)
	defer func() {
		// clean data after test
		os.RemoveAll(src)
		os.RemoveAll(dest)
	}()
	if err != nil {
		t.Fatalf("Cannot copy dir: %s", err)
	}

	bytes, err := ioutil.ReadFile(dest + filePath)
	if err != nil {
		t.Fatalf("Сannot read dest file: %s", err)
	}

	if expectedText != string(bytes) {
		t.Fatal("Expected file content does not equal real")
	}
}
