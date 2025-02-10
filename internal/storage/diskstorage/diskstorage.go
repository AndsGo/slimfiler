package diskstorage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slimfiler/internal/utils/fileutil"
	"slimfiler/internal/utils/md5"
)

// Storage is an implementation of httpcache.Storage that supplements the in-memory map with persistent storage
type Storage struct {
	DiskPath  string
	ServerURL string
}

// Get returns the response corresponding to key if present
func (c *Storage) Get(key string) (resp []byte, ETag string, err error) {
	filefullPath := path.Join(c.DiskPath, key)
	// 判断目录是否存在
	if !fileutil.IsExist(filefullPath) {
		return nil, "", os.ErrNotExist
	}
	// 获取文件目录,去除最后一个/，就是目录
	// 打开文件
	file, err := os.Open(filefullPath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	// 读取文件
	file.Read(resp)
	return resp, md5.GetMD5(resp), nil
}

func (c *Storage) GetStream(key string) (r io.ReadCloser, ETag string, err error) {
	filefullPath := path.Join(c.DiskPath, key)
	// 判断目录是否存在
	if !fileutil.IsExist(filefullPath) {
		return nil, "", os.ErrNotExist
	}
	// 获取文件目录,去除最后一个/，就是目录
	// 打开文件
	file, err := os.Open(filefullPath)
	if err != nil {
		return nil, "", err
	}
	return file, "", nil
}

// Put saves a response to the cache as key
func (c *Storage) Put(key string, resp []byte) (ETag string, err error) {
	filefullPath := path.Join(c.DiskPath, key)
	dir := filepath.Dir(filefullPath)
	fileutil.CreateDir(dir)
	if err := os.WriteFile(filefullPath, resp, 0644); err != nil {
		return "", err
	}
	return md5.GetMD5(resp), nil
}

func (c *Storage) PutStream(key string, r io.ReadCloser) (ETag string, err error) {
	filefullPath := fmt.Sprintf("%s%s", c.DiskPath, key)
	dir := filepath.Dir(filefullPath)
	// 判断目录是否存在
	if !fileutil.IsExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}
	file, err := os.Create(filefullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, r); err != nil {
		return "", err
	}
	return "", nil
}

// Delete removes the response with key from the cache
func (c *Storage) Delete(key string) error {
	return os.Remove(path.Join(c.DiskPath, key))
}
func (c *Storage) HeadObject(key string) (http.Header, error) {
	return nil, nil
}

// New returns a new Cache that will store files in basePath
func New(basePath string) *Storage {
	return &Storage{DiskPath: basePath}
}
