package utils

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

func FindFilePath(point ...string) (searched string, err error) {
	partial := path.Join(point...)
	searched, err = filepath.Abs(partial)
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(searched); errors.Is(err, os.ErrNotExist) {
		for i := 0; i <= 2; i++ {
			if filepath.Dir(searched) == "" {
				break
			}
			searched = filepath.Dir(searched)
		}
		searched = path.Join(searched, partial)
	}
	return searched, nil
}
