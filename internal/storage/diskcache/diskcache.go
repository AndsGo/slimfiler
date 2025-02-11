// Package diskcache provides an implementation of httpcache.Cache that uses the diskv package
// to supplement an in-memory map with persistent storage
package diskcache

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"slimfiler/internal/utils/fileutil"
	"slimfiler/internal/utils/md5util"

	"github.com/peterbourgon/diskv"
)

// Cache is an implementation of httpcache.Cache that supplements the in-memory map with persistent storage
type Cache struct {
	d *diskv.Diskv
}

// Get returns the response corresponding to key if present
func (c *Cache) Get(key string) (resp []byte, ETag string, err error) {
	key = keyToFilename(key)
	resp, err = c.d.Read(key)
	if err != nil {
		return []byte{}, "", err
	}
	return resp, md5util.GetMD5(resp), nil
}

func (c *Cache) GetStream(key string) (r io.ReadCloser, ETag string, err error) {
	key = keyToFilename(key)
	r, err = c.d.ReadStream(key, true)
	if err != nil {
		return nil, "", err
	}
	return r, "", nil
}

// Put saves a response to the cache as key
func (c *Cache) Put(key string, resp []byte) (ETag string, err error) {
	key = keyToFilename(key)
	err = c.d.WriteStream(key, bytes.NewReader(resp), true)
	if err != nil {
		return "", err
	}
	return md5util.GetMD5(resp), nil
}

func (c *Cache) PutStream(key string, r io.ReadCloser) (ETag string, err error) {
	key = keyToFilename(key)
	err = c.d.WriteStream(key, r, true)
	if err != nil {
		return "", err
	}
	return "", nil
}

// Delete removes the response with key from the cache
func (c *Cache) Delete(key string) error {
	key = keyToFilename(key)
	return c.d.Erase(key)
}

func (c *Cache) HeadObject(key string) (http.Header, error) {
	return nil, nil
}

func keyToFilename(key string) string {
	return md5util.GetMD5([]byte(key))
}

// New returns a new Cache that will store files in basePath
func New(basePath string) *Cache {
	dir := filepath.Dir(basePath)
	fileutil.CreateDir(dir)
	return &Cache{
		d: diskv.New(diskv.Options{
			BasePath:     basePath,
			CacheSizeMax: 100 * 1024 * 1024, // 100MB
			// For file "c0ffee", store file as "c0/ff/c0ffee"
			Transform: func(s string) []string { return []string{s[0:2], s[2:4]} },
		}),
	}
}

// NewWithDiskv returns a new Cache using the provided Diskv as underlying
// storage.
func NewWithDiskv(d *diskv.Diskv) *Cache {
	return &Cache{d}
}
