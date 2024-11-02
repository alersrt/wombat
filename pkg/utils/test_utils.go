package utils

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// FindFilePath Upward search for path described by matches
func FindFilePath(matches ...string) (searched string, err error) {
	partial := path.Join(matches...)
	searched, err = filepath.Abs(partial)
	if err != nil {
		return "", err
	}

	searched = recursiveSearch(searched, partial, len(matches))

	return searched, nil
}

func recursiveSearch(searched string, partial string, counter int) string {
	_, err := os.Stat(searched)
	if errors.Is(err, os.ErrNotExist) {
		for i := 0; i <= 2; i++ {
			tmp := filepath.Dir(searched)
			if tmp == "." || tmp == "/" {
				break
			}
			searched = tmp
		}
		searched = recursiveSearch(path.Join(searched, partial), partial, counter)
	}
	return searched
}
