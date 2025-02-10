package s3storage

import (
	"bytes"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Options aws s3服务应用层客户端
type Options struct {
	SecretId         string
	SecretKey        string
	Region           string
	Bucket           string
	Endpoint         string
	Token            string
	DisableSSL       bool
	S3ForcePathStyle bool
}

type awsS3 struct {
	Client *s3.S3
	Bucket string
}

// NewAwsS3 创建aws s3实例
func NewAwsS3(options Options) *awsS3 {
	var awsS3Instance awsS3
	awsS3Instance.Bucket = options.Bucket
	config := aws.NewConfig().WithRegion(options.Region)
	if options.DisableSSL {
		config.WithDisableSSL(true)
	}
	config.WithCredentials(credentials.NewStaticCredentials(options.SecretId, options.SecretKey, options.Token))
	config.WithEndpoint(options.Endpoint)
	if options.S3ForcePathStyle {
		config.WithS3ForcePathStyle(true)
	}
	// config := &aws.Config{
	// 	Credentials:      credentials.NewStaticCredentials(options.SecretId, options.SecretKey, options.Token),
	// 	Region:           aws.String(options.Region),
	// 	Endpoint:         aws.String(options.Endpoint),
	// 	DisableSSL:       aws.Bool(options.DisableSSL),
	// 	S3ForcePathStyle: aws.Bool(options.S3ForcePathStyle),
	// }
	sess := session.Must(session.NewSession(config))
	awsS3Instance.Client = s3.New(sess)
	return &awsS3Instance
}

// PutObject 根据内容上传文件对象
func (a *awsS3) Put(awsPath string, content []byte) (string, error) {
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
		Body:   aws.ReadSeekCloser(bytes.NewReader(content)),
	}
	resp, err := a.Client.PutObject(putObjectInput)
	if err != nil {
		return "", err
	}
	return *(resp.ETag), nil
}
func (a *awsS3) PutStream(awsPath string, r io.ReadCloser) (ETag string, err error) {
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
		Body:   aws.ReadSeekCloser(r),
	}
	resp, err := a.Client.PutObject(putObjectInput)
	if err != nil {
		return "", err
	}
	return *(resp.ETag), nil
}

// GetObject 下载文件对象内容
func (a *awsS3) Get(awsPath string) ([]byte, string, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
	}
	resp, err := a.Client.GetObject(getObjectInput)
	if err != nil {
		return nil, "", err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return content, *(resp.ETag), nil
}

func (a *awsS3) GetStream(awsPath string) (r io.ReadCloser, ETag string, err error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
	}
	resp, err := a.Client.GetObject(getObjectInput)
	if err != nil {
		return nil, "", err
	}
	return resp.Body, *(resp.ETag), nil
}

// DeleteObject 删除文件对象
func (a *awsS3) Delete(awsPath string) error {
	deleteObject := &s3.DeleteObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
	}
	_, err := a.Client.DeleteObject(deleteObject)
	if err != nil {
		return err
	}
	return nil
}

// HeadObject 获取对象元数据信息，包括md5和上次修改时间
func (a *awsS3) HeadObject(awsPath string) (string, *time.Time, error) {
	headObject := &s3.HeadObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(awsPath),
	}
	resp, err := a.Client.HeadObject(headObject)
	if err != nil {
		return "", nil, err
	}
	return *(resp.ETag), resp.LastModified, nil
}
