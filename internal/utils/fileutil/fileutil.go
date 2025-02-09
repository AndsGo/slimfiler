package fileutil

import (
	"errors"
	"net/http"
	"os"
	"strings"
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
func CreateDir(absPath string) error {
	// return os.MkdirAll(path.Dir(absPath), os.ModePerm)
	return os.MkdirAll(absPath, os.ModePerm)
}

// 判断url上是否有download参数
func SetDownload(w http.ResponseWriter, r *http.Request, fileName string) bool {
	if r.URL.Query().Get("download") == "1" { // 如果url上带有download=1参数，则认为是下载请求
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		return true
	}
	if strings.Contains(r.URL.Path, "&download=1") || strings.Contains(r.URL.Path, "?download=1") { // 如果url上带有download参数，则认为是下载请求
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		return true
	}
	w.Header().Set("Content-Disposition", "inline; filename="+fileName)
	return false
}
