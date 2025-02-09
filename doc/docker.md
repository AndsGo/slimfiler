## 构建docker

```bash
docker build -t fileserver:latest .
```
## 运行docker
```bash
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data fileserver:latest
```
## 保存镜像到文件
```bash
docker save -o fileserver.tar fileserver:latest
```
## 从文件加载镜像
```bash
docker load -i fileserver.tar
docker run -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data fileserver:latest
```

## 移除镜像
```bash
docker rmi fileserver:latest
```

## 后台运行
```bash
docker run -d --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data fileserver:latest
```
