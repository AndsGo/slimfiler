package svc

import (
	"slimfiler/internal/config"
	"slimfiler/internal/data"
	"slimfiler/internal/storage"
	"slimfiler/internal/storage/diskcache"
	"slimfiler/internal/storage/diskstorage"
	"slimfiler/internal/storage/s3storage"
	"slimfiler/internal/utils/fileutil"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type ServiceContext struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Cache   storage.Storage
	Storage storage.Storage
	Db      *data.Store
}

func NewServiceContext(c *config.Config) *ServiceContext {
	// 创建日志目录
	fileutil.CreateDir(c.Log.Path)
	// 配置 lumberjack
	logger := &lumberjack.Logger{
		Filename:   c.Log.Path + "/" + c.Log.FileName,
		MaxSize:    10,             // 最大日志文件大小（MB）
		MaxBackups: 3,              // 最多保留 3 个备份
		MaxAge:     c.Log.KeepDays, // 最多保留 7 天的日志
		Compress:   c.Log.Compress, // 是否压缩日志
	}
	svcContxet := &ServiceContext{
		Config: c,
		Logger: &logrus.Logger{Out: logger, Formatter: &logrus.JSONFormatter{}, Level: logrus.InfoLevel},
	}
	// 初始化缓存
	switch c.ImageCacheConf.Node {
	case config.DiskNode:
		svcContxet.Cache = diskcache.New(c.ImageCacheConf.DiskOptions.DiskPath)
	case config.S3Node:
		svcContxet.Cache = s3storage.NewAwsS3(s3storage.Options(c.ImageCacheConf.S3Options))
	default:
		svcContxet.Cache = storage.NopStorage
	}
	// 初始化file存储
	switch c.UploadConf.Node {
	case config.DiskNode:
		svcContxet.Storage = diskstorage.New(c.UploadConf.DiskOptions.DiskPath)
	case config.S3Node:
		svcContxet.Storage = s3storage.NewAwsS3(s3storage.Options(c.UploadConf.S3Options))
	default:
		svcContxet.Storage = storage.NopStorage
	}
	// 初始化db
	db, err := data.NewStore(data.Options{
		Path:       c.Db.Path,
		BucketName: c.Db.BucketName,
	})
	if err != nil {
		svcContxet.Logger.Errorf("db init error %s", err.Error())
		panic(err)
	}
	svcContxet.Db = db
	return svcContxet
}
