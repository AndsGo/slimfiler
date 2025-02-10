package s3storage

import (
	"crypto/md5"
	"fmt"
	"os"
	"testing"
)

func Test_Delete(t *testing.T) {
	/**
	SecretId  = "Cf98lG6JXZ0TVfste6Qt"
	SecretKey = "UkcYAAXs6XMABr59IBADfbkhj2LRfD9zsdkWEB8O"
	Region    = "CN"
	Bucket    = "test"
	Token     = DefaultToken
	Endpoint  = "http://10.0.0.119:10000"
	*/
	s3 := NewAwsS3(Options{
		SecretId:         "Cf98lG6JXZ0TVfste6Qt",
		SecretKey:        "UkcYAAXs6XMABr59IBADfbkhj2LRfD9zsdkWEB8O",
		Region:           "CN",
		Bucket:           "test",
		Token:            "",
		Endpoint:         "http://10.0.0.119:10000",
		DisableSSL:       false,
		S3ForcePathStyle: true,
	})
	err := s3.Delete("/WordPress-agency-in-Hong-Kong.jpg")
	if err != nil {
		t.Error(err)
	}
}

func Test_Put(t *testing.T) {
	/**
	SecretId  = "Cf98lG6JXZ0TVfste6Qt"
	SecretKey = "UkcYAAXs6XMABr59IBADfbkhj2LRfD9zsdkWEB8O"
	Region    = "CN"
	Bucket    = "test"
	Token     = DefaultToken
	Endpoint  = "http://10.0.0.119:10000"
	*/
	s3 := NewAwsS3(Options{
		SecretId:         "Cf98lG6JXZ0TVfste6Qt",
		SecretKey:        "UkcYAAXs6XMABr59IBADfbkhj2LRfD9zsdkWEB8O",
		Region:           "CN",
		Bucket:           "test",
		Token:            "",
		Endpoint:         "http://10.0.0.119:10000",
		DisableSSL:       false,
		S3ForcePathStyle: true,
	})
	localFilePath := "C:\\Users\\Administrator\\Downloads\\WordPress-agency-in-Hong-Kong.jpg"

	t.Logf("local file path %s", localFilePath)
	fileContent, err := os.ReadFile(localFilePath)
	if err != nil {
		t.Fatalf("read file error: %s!", err.Error())
		return
	}
	// 测试时修改aws路径
	awsPath := "/WordPress-agency-in-Hong-Kong.jpg"
	hash := md5.New()
	hash.Write(fileContent)
	fmt.Printf("md5:%x", hash.Sum(nil))
	tag, err := s3.Put(awsPath, fileContent)
	if err != nil {
		t.Fatalf("put file error: %s!", err.Error())
		return
	}
	t.Logf("aws path %s, tag %s", awsPath, tag)
}

func Test_Get(t *testing.T) {
	s3 := NewAwsS3(Options{
		SecretId:         "Cf98lG6JXZ0TVfste6Qt",
		SecretKey:        "UkcYAAXs6XMABr59IBADfbkhj2LRfD9zsdkWEB8O",
		Region:           "CN",
		Bucket:           "test",
		Token:            "",
		Endpoint:         "http://10.0.0.119:10000",
		DisableSSL:       false,
		S3ForcePathStyle: true,
	})
	awsPath := "/WordPress-agency-in-Hong-Kong.jpg"
	data, tag, err := s3.Get(awsPath)
	if err != nil {
		t.Fatalf("get file error: %s!", err.Error())
		return
	}
	t.Logf("aws path %s, tag %s", awsPath, tag)
	t.Logf("data md5:%x", md5.Sum(data))
}
