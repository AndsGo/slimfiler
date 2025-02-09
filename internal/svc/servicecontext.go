package svc

import (
	"fileserver/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type ServiceContext struct {
	Config *config.Config
	Logger *logrus.Logger
}

func NewServiceContext(c *config.Config) *ServiceContext {
	// 配置 lumberjack
	logger := &lumberjack.Logger{
		Filename:   c.Log.Path + "/" + c.Log.FileName,
		MaxSize:    10,             // 最大日志文件大小（MB）
		MaxBackups: 3,              // 最多保留 3 个备份
		MaxAge:     c.Log.KeepDays, // 最多保留 7 天的日志
		Compress:   c.Log.Compress, // 是否压缩日志
	}
	return &ServiceContext{
		Config: c,
		Logger: &logrus.Logger{Out: logger, Formatter: &logrus.JSONFormatter{}, Level: logrus.InfoLevel},
	}
}
