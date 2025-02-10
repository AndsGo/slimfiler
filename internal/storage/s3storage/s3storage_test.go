package s3storage

import (
	"crypto/md5"
	"fmt"
	"net/http"
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
	awsPath := "WordPress-agency-in-Hong-Kong.jpg"
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

func Test_Put2(t *testing.T) {
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
	fileContent, err := os.Open(localFilePath)
	if err != nil {
		t.Fatalf("read file error: %s!", err.Error())
		return
	}
	// 测试时修改aws路径
	awsPath := "WordPress-agency-in-Hong-Kong.jpg"
	tag, err := s3.PutStream(awsPath, fileContent)
	if err != nil {
		t.Fatalf("put file error: %s!", err.Error())
		return
	}
	t.Logf("aws path %s, tag %s", awsPath, tag)
}

func Test_PutStream(t *testing.T) {
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
	response, err := http.Get("https://img.kakaclo.com/image%2FFSZW17727%2FFSZW17727_WHE_S_NUB%2F5a12dad52968410cad955e01ab893b21.jpg")
	if err != nil {
		t.Fatalf("get file error: %s!", err.Error())
	}
	// // 创建一个读写管道
	// pr, pw := io.Pipe()

	// // 异步读取 HTTP 响应并写入管道
	// go func() {
	// 	defer pw.Close()
	// 	if _, err := io.Copy(pw, response.Body); err != nil {
	// 		fmt.Println("复制响应体错误:", err)
	// 	}
	// }()
	// 测试时修改aws路径
	tag, err := s3.PutStream("https://img.kakaclo.com/image%2FFSZW17727%2FFSZW17727_WHE_S_NUB%2F5a12dad52968410cad955e01ab893b21.jpg", response.Body)
	if err != nil {
		t.Fatalf("put file error: %s!", err.Error())
		return
	}
	t.Logf("aws path %s, tag %s", "awsPath", tag)
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
