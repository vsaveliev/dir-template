package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// tarGzArchivator creates *.tar.gz archives from files, dirs
type tarGzArchivator struct {
	Src  string
	Dest string

	tw *tar.Writer
}

// NewTarGzArchivator returns a new instance of TarGzArchivator
func newTarGzArchivator(src, dest string) *tarGzArchivator {
	return &tarGzArchivator{
		Src:  src,
		Dest: dest,
	}
}

// CreateTarGzArchive creates dest.tar.gz archive
func CreateTarGzArchive(src, dest string) error {
	t := newTarGzArchivator(src, dest)
	return t.create()
}

// tarGzFile adds file in tar gz writer
func (t *tarGzArchivator) tarGzFile(relativePath string, fi os.FileInfo) error {
	fr, err := os.Open(t.Src + "/" + relativePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", err)
	}
	defer fr.Close()

	h := new(tar.Header)
	h.Name = relativePath
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = t.tw.WriteHeader(h)
	if err != nil {
		return fmt.Errorf("cannot write header to archive: %s", err)
	}

	_, err = io.Copy(t.tw, fr)
	if err != nil {
		return fmt.Errorf("cannot to copy file in archive: %s", err)
	}

	return nil
}

// tarGzDir adds dir in tar gz writer
func (t *tarGzArchivator) tarGzDir(relativePath string) error {
	dir, err := os.Open(t.Src + "/" + relativePath)
	if err != nil {
		return fmt.Errorf("cannot open src dir: %s", err)
	}
	defer dir.Close()

	objects, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("cannot read src dir: %s", err)
	}

	for _, obj := range objects {
		srcObjPath := relativePath + "/" + obj.Name()

		if obj.IsDir() {
			err = t.tarGzDir(srcObjPath)
			if err != nil {
				return err
			}
		} else {
			err = t.tarGzFile(srcObjPath, obj)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Create makes *.tar.gz
func (t *tarGzArchivator) create() error {
	info, err := os.Stat(t.Src)
	if err != nil {
		return fmt.Errorf("cannot get info of src dir: %s", err)
	}

	fw, err := os.Create(t.Dest)
	if err != nil {
		return fmt.Errorf("cannot create archive: %s", err)
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	t.tw = tar.NewWriter(gw)
	defer t.tw.Close()

	if info.IsDir() {
		return t.tarGzDir("")
	}

	// srcPath = dir of src file
	t.Src = strings.TrimSuffix(t.Src, "/"+info.Name())

	return t.tarGzFile(info.Name(), info)
}
