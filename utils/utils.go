// Package utils contains utility functions that are used everywhere in the app.
package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// StringInSlice check wether a string is contain in a string slice.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IndexOfString returns the index of a string in a string slice
func IndexOfString(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return 0
}

// FileName returns the name of a file, without the path and extension.
func FileName(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name
}

// DatedName returns the name of a file with a date appended, without extension.
// It is used for savestates and screenshot names.
func DatedName(path string) string {
	name := FileName(path)
	date := time.Now().Format("2006-01-02-15-04-05")
	return name + "@" + date
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

// CaptureOutput executes a function and capture all the text outputted to logs.
// Used in unit tests.
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

// AllFilesIn recursively builds a list of the files in a given directory
func AllFilesIn(dir string) ([]string, error) {
	file := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			file = append(file, path)
		}
		return nil
	})
	return file, err
}

// CoreExt returns the libretro core extension for the current OS
func CoreExt() string {
	exts := map[string]string{
		"darwin":  ".dylib",
		"linux":   ".so",
		"windows": ".dll",
	}

	return exts[runtime.GOOS]
}

// LinesInFile counts the number of lines in a file
func LinesInFile(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
