package plexrenamer

import (
	"os"
	"path/filepath"
	"strings"
)

type FileResult struct {
	From string
	To   string
}

func ScanDir(fromDir, toDir string, dryrun bool) ([]FileResult, error) {
	results := make([]FileResult, 0)

	c, err := os.ReadDir(fromDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range c {
		if entry.IsDir() {
			continue
		}

		oldFilePath := filepath.Join(fromDir, entry.Name())
		newFilePath, err := GetNewFilePath(oldFilePath)
		if err != nil {
			continue
		}
		newFilePath = filepath.Join(toDir, newFilePath)

		if !dryrun {
			err = os.MkdirAll(filepath.Dir(newFilePath), os.ModePerm)
			if err != nil {
				return results, err
			}

			err = os.Rename(oldFilePath, newFilePath)
			if err != nil {
				return results, err
			}
		}

		results = append(results, FileResult{
			From: oldFilePath,
			To:   newFilePath,
		})
	}

	return results, nil
}

func GetNewFilePath(oldFilePath string) (string, error) {
	filename := filepath.Base(oldFilePath)
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext)

	pi, err := Parse(filename)
	if err != nil {
		return "", err
	}

	newFilePath, err := PlexFormat(pi, ext)
	if err != nil {
		return "", err
	}

	return newFilePath, nil
}
