## 业务场景
这个项目主要用于搭建公司内部的文件管理系统。加快访问图片的速度。因此需要考虑以下问题：
1. 文件上传下载
2. 文件代理/缓存
3. 图片压缩/缓存
适合的页面需要查看大量图片和文件频繁访问的业务场景：
1. 商品系统
2. ai生成图
3. 内网cdn (文件代理)
## 开发思路
开发一个文件服务器，提供文件上传下载，代理，以及图片压缩等功能。其中图片压缩功能依赖[imageprocess](https://github.com/AndsGo/imageprocess)。

需要包含以下功能：
1. 文件上传下载
2. 文件代理
3. 图片压缩

## 1. 文件上传下载
基于 `http` 库
文件上传下载功能使用，使用multipart/form-data上传文件，使用uuid作为文件名，使用`http`包的静态文件服务来提供文件下载。
> 需要考虑以下问题：
1. 文件上传大小限制
2. 文件上传目录
3. 文件上传后缀限制
4. 文件上传重名处理
5. 文件重复秒传
> 需要对文件存储支持:
1. 本地存储
2. 云存储 (如：阿里云OSS,s3)

## 2. 文件代理
基于 `http` 库，参考[imageproxy](https://github.com/willnorris/imageproxy)实现
> 需要考虑以下问题：
1. 代理请求头
2. 代理请求参数
3. 文件代理缓存
> 需要对文件缓存支持：
1. 本地缓存
2. 云存储 (如：阿里云OSS,s3)

## 3. 图片压缩
基于 `imageprocess` 库，参考[imageprocess](https://github.com/AndsGo/imageprocess)实现

## 4. 依赖库
1. [diskv](https://github.com/peterbourgon/diskv) 本地缓存
2. [imageprocess](https://github.com/AndsGo/imageprocess) 图片压缩，水印，格式转换
3. [s3](https://github.com/aws/aws-sdk-go) AWS S3 可以连接支持 s3协议的文件服务器
4. [bbolt](https://github.com/etcd-io/bbolt) 本地数据库，数据库存储文件信息
5. [yaml](https://github.com/go-yaml/yaml/tree/v3.0.1) yaml解析
6. [logrus](https://github.com/sirupsen/logrus) 日志
7. [lumberjack](https://github.com/natefinch/lumberjack) 日志切割
