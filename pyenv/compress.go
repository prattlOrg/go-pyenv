package pyenv

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ZIP_FILE_EXT = ".zip"

// Stolen from https://gosamples.dev/zip-file/
func compressDir(source, target string) error {
	log.Printf("Compressing %v into %v\n", source, target)
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

// Stolen from https://gosamples.dev/unzip-file/
func unzipSource(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	// log.Printf("unzipping %v into %v\n", source, destination)
	if err != nil {
		return err
	}

	sourceDir, hasSuffix := strings.CutSuffix(source, ZIP_FILE_EXT)
	if !hasSuffix {
		log.Fatalf("Expected source file to end in .zip, got = %v\n", source)
	}

	sourceDir = filepath.Base(sourceDir) + string(os.PathSeparator)

	for _, f := range reader.File {
		// log.Printf("unzipping file: %v\nwith source dir: %v\n", f.Name, sourceDir)
		err := unzipFile(f, destination, sourceDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination, sourceDir string) error {
	if f.Name == sourceDir {
		return nil
	}

	fileName := strings.TrimPrefix(f.Name, sourceDir)
	filePath := filepath.Join(destination, fileName)
	os.Mkdir(destination, 0o777)

	// log.Printf("Unzipping file: %v\n", filePath)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
