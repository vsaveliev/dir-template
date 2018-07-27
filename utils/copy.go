package utils

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies file from srcFilePath to destFilePath with current filemode
func CopyFile(srcFilePath, destFilePath string) (err error) {
	if srcFilePath == "" {
		return fmt.Errorf("src file path is empty")
	}

	if destFilePath == "" {
		return fmt.Errorf("dest file path is empty")
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("cannot open src file: %s", err)
	}
	defer srcFile.Close()

	srcFileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return fmt.Errorf("cannot get info of src file: %s", err)
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("cannot create dest file: %s", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("cannot copy file: %s", err)
	}

	err = os.Chmod(destFilePath, srcFileInfo.Mode())
	if err != nil {
		os.Remove(destFilePath)
		return fmt.Errorf("cannot copy file mode: %s", err)
	}

	return nil
}

// CopyDir copies dir from srcDirPath to destDirPath with current filemodes
func CopyDir(srcDirPath, destDirPath string) error {
	if srcDirPath == "" {
		return fmt.Errorf("src path is empty")
	}

	if destDirPath == "" {
		return fmt.Errorf("dest path is empty")
	}

	dirInfo, err := os.Stat(srcDirPath)
	if err != nil {
		return fmt.Errorf("cannot get info of src dir: %s", err)
	}

	srcDir, err := os.Open(srcDirPath)
	if err != nil {
		return fmt.Errorf("cannot open src dir: %s", err)
	}
	defer srcDir.Close()

	// objects = files or child directories
	dirObjects, err := srcDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("cannot read src dir: %s", err)
	}

	err = os.MkdirAll(destDirPath, dirInfo.Mode())
	if err != nil {
		return fmt.Errorf("cannot create dest dir %s: %s", destDirPath, err)
	}

	for _, obj := range dirObjects {
		srcObjPath := srcDirPath + "/" + obj.Name()
		destObjPath := destDirPath + "/" + obj.Name()

		if !obj.IsDir() {
			err = CopyFile(srcObjPath, destObjPath)
			if err != nil {
				os.Remove(destDirPath)
				return err
			}

			continue
		}

		err = CopyDir(srcObjPath, destObjPath)
		if err != nil {
			os.Remove(destDirPath)
			return err
		}
	}

	return nil
}
