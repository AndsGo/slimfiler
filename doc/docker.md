## 构建docker

```bash
docker build -t slimfiler:latest .
```
## 运行docker
```bash
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data slimfiler:latest
```
## 保存镜像到文件
```bash
docker save -o slimfiler.tar slimfiler:latest
```
## 从文件加载镜像
```bash
docker load -i slimfiler.tar
docker run -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data slimfiler:latest
```

## 移除镜像
```bash
docker rmi slimfiler:latest
```

## 后台运行
```bash
docker run -d --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/data>:/data slimfiler:latest
```
