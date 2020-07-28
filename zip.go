package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 压缩文件
// src: 待压缩文件外层目录
// dst: 压缩文件目录
func Zip(src, dst string) error {
	// Get a Buffer to Write To
	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)
	// Add some files to the archive.
	err = addFiles(w, src, "")

	if err != nil {
		return err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(filepath.Join(basePath, file.Name()))
		if !file.IsDir() {
			// file
			p := filepath.Join(basePath, file.Name())
			dat, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}

			p = filepath.Join(baseInZip, file.Name())
			// Add some files to the archive.
			f, err := w.Create(p)
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Dir
			newBase := filepath.Join(basePath, file.Name())
			err = addFiles(w, newBase, filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 解压文件
// src: 待解压文件，“/User/Desktop/somefile.zip”
// dst: 解压后文件目录，“/User/Desktop/dst”
func Unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}
