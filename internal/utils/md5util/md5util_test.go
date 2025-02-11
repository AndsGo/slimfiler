package md5util

import "testing"

func Test_GetMD5(t *testing.T) {
	t.Log(GetMD5([]byte("test")))
	// 098f6bcd4621d373cade4e832627b4f6
}
