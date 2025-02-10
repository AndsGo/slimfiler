package config

type Config struct {
	Name string `yaml:"Name"`
	Port int    `yaml:"Port"`
	Log  struct {
		FileName string `yaml:"FileName"`
		Mode     string `yaml:"Mode"`
		Path     string `yaml:"Path"`
		Level    string `yaml:"Level"`
		Compress bool   `yaml:"Compress"`
		KeepDays int    `yaml:"KeepDays"`
	} `yaml:"Log"`

	UploadConf struct {
		MaxImageSize int64       `yaml:"MaxImageSize"`
		MaxVideoSize int64       `yaml:"MaxVideoSize"`
		MaxAudioSize int64       `yaml:"MaxAudioSize"`
		MaxOtherSize int64       `yaml:"MaxOtherSize"`
		ServerURL    string      `yaml:"ServerURL"`
		Node         Node        `yaml:"Node"`
		DiskOptions  DiskOptions `yaml:"DiskOptions"`
		S3Options    S3Options   `yaml:"S3Options"`
	} `yaml:"UploadConf"`
	ImageCacheConf struct {
		Node        Node        `yaml:"Node"`
		DiskOptions DiskOptions `yaml:"DiskOptions"`
		S3Options   S3Options   `yaml:"S3Options"`
	} `yaml:"ImageCacheConf"`
	Db struct {
		Path       string `yaml:"Path"`
		BucketName string `yaml:"BucketName"`
	} `yaml:"Db"`
}

type DiskOptions struct {
	DiskPath  string `yaml:"DiskPath"`
	ServerURL string `yaml:"ServerURL"`
}

type S3Options struct {
	SecretId         string `yaml:"SecretId"`
	SecretKey        string `yaml:"SecretKey"`
	Region           string `yaml:"Region"`
	Bucket           string `yaml:"Bucket"`
	Endpoint         string `yaml:"Endpoint"`
	Token            string `yaml:"Token"`
	DisableSSL       bool   `yaml:"DisableSSL"`
	S3ForcePathStyle bool   `yaml:"S3ForcePathStyle"`
}

type Node string

const (
	DiskNode Node = "DiskOptions"
	S3Node   Node = "S3Options"
)
