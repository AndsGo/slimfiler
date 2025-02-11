package fileutil

import (
	"errors"
	"os"
)

// IsExist checks if a file or directory exists.
// Play: https://go.dev/play/p/nKKXt8ZQbmh
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// CreateFile create a file in path.
// Play: https://go.dev/play/p/lDt8PEsTNKI
func CreateFile(path string) bool {
	file, err := os.Create(path)
	if err != nil {
		return false
	}

	defer file.Close()
	return true
}

// CreateDir create directory in absolute path. param `absPath` like /a/, /a/b/.
// Play: https://go.dev/play/p/qUuCe1OGQnM
func CreateDir(dir string) error {
	// 判断目录是否存在
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
