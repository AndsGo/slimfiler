## 介绍
这是一个Go语言开发的文件服务器，提供文件上传下载，代理，以及图片压缩等功能。其中图片处理功能使用兼容`aliyun oss` 参数。

## 安装
linux [下载](https://github.com/AndsGo/SlimFiler/releases/download/v1.1.0/slimfiler)

windows [下载](https://github.com/AndsGo/SlimFiler/releases/download/v1.1.0/slimfiler.exe)

配置文件 [下载](https://github.com/AndsGo/slimfiler/releases/download/v1.1.0/conf.yaml)

> 运行
```
slimfiler -f ./conf.yaml
```
or
```
slimfiler.exe -f ./conf.yaml
```
配置文件具体参数参考[conf.md](./doc/conf.md)
## 源码构建
在根目录下执行
```
go build 
```
## 运行
构建完成后执行
```
slimfiler.exe -f ./conf/conf.yaml
```
如果是linux系统，则执行
```
./slimfiler -f ./conf/conf.yaml
```
## 使用方法
### 1.上传文件
```bash
curl --location --request POST '127.0.0.1:8000/upload' \
--form 'file=@"/C:/Users/Administrator/OneDrive/Desktop/1c7b79352c7679274f242d72cec433b2.jpg"' \
--form 'directory="data/test"'
```
response:
```json
{
    "code": 10000,
    "data": {
        "name": "1c7b79352c7679274f242d72cec433b2.jpg",
        "url": "/data/test/2025-02-07/987cb8c4-f7f8-4b1e-940a-f9dee03e6ec8.jpg"
    },
    "msg": ""
}
```
> 支持form参数:
1. `file` (必填)
2. `directory` (选填)
3. `md5` (选填) 用于验证文件是否已存在，实现秒传
### 2.查看文件
使用上传文件的返回的url访问文件
```bash
curl --location --request GET 'http://127.0.0.1:8080/data/test/2025-02-07/987cb8c4-f7f8-4b1e-940a-f9dee03e6ec8.jpg'
```
增加图片处理:
http://127.0.0.1:8000/data/test/2025-02-07/987cb8c4-f7f8-4b1e-940a-f9dee03e6ec8.jpg?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp
### 3.代理
代理部分功能参考[imageproxy](https://github.dev/willnorris/imageproxy)实现
```bash
curl --location --request GET 'http://127.0.0.1:8000/proxy/https://raw.githubusercontent.com/AndsGo/imageprocess/refs/heads/main/examples/example.jpg'

```
增加图片处理:

http://127.0.0.1:8000/proxy/https://raw.githubusercontent.com/AndsGo/imageprocess/refs/heads/main/examples/example.jpg?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp

## 图片处理
图片处理依赖[imageprocess](https://github.com/AndsGo/imageprocess),图片查看和代理都支持图片压缩。
图片处理兼容 `aliyun oss` 的图片处理，在url后面添加参数`x-oss-process`即可。支持图片格式：`WEBP`,`JPG`,`JPEG`,`PNG`,`BMP`,`TIFF`,`GIF`。 
### 常用例子
1. 压缩图片 宽度为360，高度为540
```
?x-oss-process=image/resize,w_360,h_540,m_fill
```
2. 压缩图片 宽度为360，高度为540，并转为webp格式
```
?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp
```
3. 压缩图片 宽度为360，高度为540，并转为webp格式，并设置质量为80
```
?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp/quality,80
```
4. 压缩图片 宽度为360，高度为540，并转为webp格式，并设置质量为80，并设置旋转角度为90度
```
?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp/quality,80/rotate,90
```
5. 压缩图片 宽度为360，高度为540，并转为webp格式，并设置质量为80，并设置旋转角度为90度，并设置水印
```
?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp/quality,80/rotate,90/watermark,text_Hello
```
更多设置请查看[imageprocess](https://github.com/AndsGo/imageprocess)

## 配置文件
配置文件在conf/conf.yaml，你在实际使用时应该进行修改。
### 示例
```yaml
Name: aigc
Port: 8000
UploadConf:
  MaxImageSize: 33554432  # 32 mb
  MaxVideoSize: 1073741824 # 1gb
  MaxAudioSize: 33554432  # 32mb
  MaxOtherSize: 10485760  # 10 mb
  ServerURL: "" # nginx path
  Node: S3Options
  DiskOptions:
    DiskPath: "./data/file/public"
  S3Options:
    Bucket: "aigc"
    Endpoint: "http://10.0.0.119:10000"
    SecretId: "xx"
    SecretKey: "xx"
    Region: "CN"
    Token: ""
    S3ForcePathStyle: true
    DisableSSL: false
PorxyCacheConf:
  Node: S3Options
  DiskOptions:
    DiskPath: "/tmp/cache" # 一般在 /tmp 下 方便自动清理
  S3Options:
    Bucket: "test"
    Endpoint: "http://10.0.0.119:10000"
    SecretId: "xx"
    SecretKey: "xx"
    Region: "CN"
    Token: ""
    S3ForcePathStyle: true
    DisableSSL: false
Db:
  Path: ./data/file/db/fileManager.db
  BuketName: "fileManager" # difualt bucket name
Log:
  FileName: fileManagerLogger
  # Mode: file # file or console
  Path: ./data/logs/file/api
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100
```
具体配置项含义和使用方法请查看[conf.md](./doc/conf.md)
### 常用配置
1.[上传文件访问地址配置](./doc/conf.md#UploadConf.ServerURL)

2.[上传文件保存配置](./doc/conf.md#UploadConf.Node) 支持本地和s3协议

2.[文件proxy缓存配置](./doc/conf.md#ImageCacheConf.Node) 支持本地和s3协议

## Docker build
Docker build 查看 [docker.md](./doc/docker.md)

## 开发思路
查看[develop.md](./doc/develop.md)