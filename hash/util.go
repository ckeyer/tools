package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ckeyer/commons/fileutils"
)

func GetFiles(path string, excludes []string) ([]string, error) {
	path = strings.TrimRight(path, "/")
	includes := []string{}

	err := filepath.Walk(path, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || err != nil {
			return err
		}

		matched, err := fileutils.Matches(strings.TrimPrefix(fpath, path+"/"), excludes)
		if err != nil || matched {
			return err
		}

		if f.Mode().IsRegular() {
			includes = append(includes, path+"/"+fpath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return includes, nil
}
