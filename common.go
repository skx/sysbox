// common.go - some routines that are used by multiple sub-commands

package main

import (
	"os"
	"path/filepath"
	"strings"
)

// FindFiles finds any file beneath the given prefix-directory which contains
// a suffix included in the list.
func FindFiles(path string, suffixes []string) ([]string, error) {

	var results []string

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		// Null info?  That probably means that the
		// destination we're trying to walk doesn't exist.
		if f == nil {
			return nil
		}

		if !f.IsDir() {
			for _, suffix := range suffixes {
				if strings.HasSuffix(path, suffix) {
					results = append(results, path)
				}
			}
		}
		return err
	})

	return results, err
}
