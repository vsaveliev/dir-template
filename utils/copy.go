package utils

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies file from src to dest (with file mode)
func CopyFile(src, dest string) (err error) {
	if src == "" {
		return fmt.Errorf("src file is empty")
	}

	if dest == "" {
		return fmt.Errorf("dest file is empty")
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open src file: %s", err)
	}
	defer srcFile.Close()

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot get info of src file: %s", err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("cannot create dest file: %s", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("cannot copy file: %s", err)
	}

	err = os.Chmod(dest, srcFileInfo.Mode())
	if err != nil {
		os.Remove(dest)
		return fmt.Errorf("cannot copy file mode: %s", err)
	}

	return nil
}

// CopyDir copies dir from src to desc (with file mode)
func CopyDir(src, dest string) error {
	if src == "" {
		return fmt.Errorf("src dir is empty")
	}

	if dest == "" {
		return fmt.Errorf("dest dir is empty")
	}

	dirInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot get info of src dir: %s", err)
	}

	srcDir, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open src dir: %s", err)
	}
	defer srcDir.Close()

	// objects = files or child directories
	dirObjects, err := srcDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("cannot read src dir: %s", err)
	}

	err = os.MkdirAll(dest, dirInfo.Mode())
	if err != nil {
		return fmt.Errorf("cannot create dest dir %s: %s", dest, err)
	}

	for _, obj := range dirObjects {
		srcObjPath := src + "/" + obj.Name()
		destObjPath := dest + "/" + obj.Name()

		if !obj.IsDir() {
			err = CopyFile(srcObjPath, destObjPath)
			if err != nil {
				os.Remove(dest)
				return err
			}

			continue
		}

		err = CopyDir(srcObjPath, destObjPath)
		if err != nil {
			os.Remove(dest)
			return err
		}
	}

	return nil
}
