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

// Filename returns the name of a file, without the path and extension.
func Filename(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

// CaptureOutput executes a function and capture all the text outputed to logs.
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
func AllFilesIn(dir string) []string {
	file := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			file = append(file, path)
		}
		return nil
	})
	return file
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
