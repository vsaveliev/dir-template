package dirtemplate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vsaveliev/dirtemplate/utils"
)

func (d DirTemplate) copyDir() error {
	_, err := os.Stat(d.Dest)
	if err == nil {
		return fmt.Errorf("dest path %s already exists", d.Dest)
	}

	err = utils.CopyDir(d.Src, d.Dest)
	if err != nil {
		return fmt.Errorf("cannot copy dir: %s", err)
	}

	// paste app name and etc. in template paths
	pathsMapForMove, err := executeFromMap(d.ReplacePaths, d.LeftDelim, d.RightDelim, d.Data)
	if err != nil {
		return fmt.Errorf("cannot process templates for rename/move paths: %s", err)
	}
	for srcPath, destPath := range pathsMapForMove {
		_, err := os.Stat(d.Dest + "/" + srcPath)
		if err != nil {
			// dir doesn't exist
			continue
		}

		err = os.Rename(d.Dest+"/"+srcPath, d.Dest+"/"+destPath)
		if err != nil {
			return fmt.Errorf("cannot rename/move path %s to %s", d.Dest+"/"+srcPath, d.Dest+"/"+destPath)
		}
	}

	return nil
}

// Process all templates from source path
func (d DirTemplate) executeDir() error {
	// check if the source dir exist
	src, err := os.Stat(d.Src)
	if err != nil {
		return err
	}

	if src.IsDir() {
		return d.processTemplatesDir(d.Src)
	}

	return d.processTemplateFile(d.Src)
}

// ProcessTemplatesDir process templates by path
func (d DirTemplate) processTemplatesDir(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open dir: %s", err)
	}
	objects, err := directory.Readdir(-1)
	if err != nil {
		return fmt.Errorf("cannot read dir: %s", err)
	}

	for _, obj := range objects {
		objectPath := path + "/" + obj.Name()
		if d.skipObject(objectPath) {
			continue
		}

		if obj.IsDir() {
			err = d.processTemplatesDir(objectPath)
			if err != nil {
				return fmt.Errorf("cannot process dir %s: %s", objectPath, err)
			}
		} else {
			err = d.processTemplateFile(objectPath)
			if err != nil {
				return fmt.Errorf("cannot process file %s: %s", objectPath, err)
			}
		}
	}

	return nil
}

// processTemplateFile process template file by path
func (d DirTemplate) processTemplateFile(path string) error {
	fileName := filepath.Base(path)

	tpl, err := template.New(fileName).Funcs(d.FuncMap).Delims(d.LeftDelim, d.RightDelim).ParseFiles(path)
	if err != nil {
		return fmt.Errorf("cannot parse file: %s", err)
	}

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("cannot create file: %s", err)
	}

	err = tpl.Execute(f, d.Data)
	if err != nil {
		return fmt.Errorf("cannot execute template: %s", err)
	}

	return nil
}

// skipObject checks by list that path can be skipped
func (d DirTemplate) skipObject(path string) bool {
	for _, skipPath := range d.SkipPaths {
		if path == strings.TrimRight(d.Src+"/"+skipPath, "/") {
			return true
		}
	}

	return false
}
