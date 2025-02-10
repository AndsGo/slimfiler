package storage

import (
	"io"
	"net/http"
)

// The Cache interface defines a cache for storing arbitrary data.  The
type Storage interface {
	// Get retrieves the Storage data for the provided key.
	Get(key string) (data []byte, ETag string, err error)

	// 流式获取
	GetStream(key string) (r io.ReadCloser, ETag string, err error)

	// Put Storage the provided data.
	Put(key string, data []byte) (ETag string, err error)

	// 流式提交
	PutStream(key string, r io.ReadCloser) (ETag string, err error)

	// Delete deletes the Storage data at the specified key.
	Delete(key string) error

	HeadObject(key string) (http.Header, error)
}

// NopStorage provides a no-op cache implementation that doesn't actually cache anything.
var NopStorage = &nopStorage{}

type nopStorage struct{}

func (c nopStorage) Get(string) (data []byte, ETag string, err error) { return nil, "", nil }

// 流式获取
func (c nopStorage) GetStream(key string) (r io.ReadCloser, ETag string, err error) {
	return nil, "", nil
}
func (c nopStorage) Put(string, []byte) (ETag string, err error) { return "", nil }

// 流式提交
func (c nopStorage) PutStream(key string, r io.ReadCloser) (ETag string, err error) {
	return "", nil
}
func (c nopStorage) Delete(string) error                    { return nil }
func (c nopStorage) HeadObject(string) (http.Header, error) { return nil, nil }
