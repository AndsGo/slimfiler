## 介绍
这是一个Go语言开发的文件服务器，提供文件上传下载，代理，以及图片压缩等功能。其中图片压缩功能依赖[imageprocess](https://github.com/AndsGo/imageprocess)。
## 打包
在根目录下执行
```
go build 
```
## 运行
```
fileserver.exe -f ./conf/conf.yaml
```
如果是linux系统，则执行
```
./fileserver -f ./conf/conf.yaml
```
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
  PublicStorePath: /data/file/public  # public path for every one access e.g. nginx path
  ServerURL: "" # nginx path
  DeleteFileWithCloud: true # whether to delete the cloud files when delete local data | 是否删除本地文件的同时删除云端文件
Log:
  FileName: fileManagerLogger
#   Mode: console 暂不支持
  Path: /data/logs/file/api
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100
```
更多配置项请查看[conf.md](./doc/conf.md)
## 使用方法
### 1.上传文件
```
curl --location --request POST '127.0.0.1:8000/upload' \
--form 'file=@"/C:/Users/Administrator/OneDrive/Desktop/1c7b79352c7679274f242d72cec433b2.jpg"' \
--form 'directory="data/test"'
```
response:
```
{
    "code": 10000,
    "data": {
        "name": "1c7b79352c7679274f242d72cec433b2.jpg",
        "url": "/data/test/2025-02-07/987cb8c4-f7f8-4b1e-940a-f9dee03e6ec8.jpg"
    },
    "msg": ""
}
```
### 2.查看文件
使用上传文件的返回的url访问文件
```
curl --location --request GET 'http://127.0.0.1:8080/data/test/2025-02-07/987cb8c4-f7f8-4b1e-940a-f9dee03e6ec8.jpg'
```
### 3.代理
```
curl --location --request GET 'http://127.0.0.1:8080/proxy/https://camo.githubusercontent.com/253a4a7e49274a9080bf1264bf001f933787b6850aa94b28442b82020bf5772b/68747470733a2f2f68656c702d7374617469632d616c6979756e2d646f632e616c6979756e63732e636f6d2f6173736574732f696d672f7a682d434e2f303331373738393636312f703532373137312e6a7067?x-oss-process=image/resize,w_360,h_540,m_fill/format,webp'
```

## 图片处理
图片处理依赖[imageprocess](https://github.com/AndsGo/imageprocess),图片查看和代理都支持图片压缩。

## Docker build
Docker build 查看 [docker.md](./doc/docker.md)

## 开发思路
查看[develop.md](./doc/develop.md)