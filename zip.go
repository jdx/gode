package gode

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func extractZip(zipfile, root string) error {
	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, f := range archive.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		path := filepath.Join(root, f.Name)
		switch {
		case f.FileInfo().IsDir():
			if err := os.Mkdir(path, f.Mode()); err != nil {
				return err
			}
		default:
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
