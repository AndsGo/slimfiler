package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slimfiler/internal/config"
	"slimfiler/internal/svc"
	"slimfiler/internal/utils/fileutil"
	"strings"

	"github.com/AndsGo/imageprocess"
	"github.com/google/uuid"
)

func ProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只支持get
		if r.Method != http.MethodGet {
			http.Error(w, "only support get", http.StatusBadRequest)
			return
		}
		logx := svcCtx.Logger
		// 获取url path
		// 获取文件名称 url,params
		parts := strings.Split(strings.TrimPrefix(r.RequestURI, "/proxy/"), "x-oss-process=")
		// url解码
		fileURL, _ := url.QueryUnescape(parts[0])
		fileURL = strings.TrimRight(fileURL, "&")
		fileURL = strings.TrimRight(fileURL, "?")
		fileURL = strings.Replace(fileURL, "http:/", "http://", -1)
		fileURL = strings.Replace(fileURL, "https:/", "https://", -1)
		// 获取 response 文件类型
		// response to Reader
		response, file, err := httpCache(fileURL, svcCtx)
		fileType := strings.Split(response.Header.Get("Content-Type"), "/")[0]
		if fileType != "image" && fileType != "video" && fileType != "audio" {
			fileType = "other"
		}
		if err != nil {
			logx.Errorf("newFunction %s", err.Error())
			http.Error(w, fmt.Sprintf("newFunction %s", err.Error()), http.StatusInternalServerError)
			return
		}
		// 设置直接下载
		fileutil.SetDownload(w, r, uuid.NewString()+"."+strings.Split(strings.Split(response.Header.Get("Content-Type"), "/")[1], "?")[0])
		if fileType != "image" {
			if _, err := io.Copy(w, file); err != nil {
				logx.Errorf("Copy file error: %s", err.Error())
				http.Error(w, "Failed to send file", http.StatusInternalServerError)
			}
			return
		}
		f, err := imageprocess.FormatFromExtension(strings.Split(strings.Split(response.Header.Get("Content-Type"), "/")[1], "?")[0])
		if err != nil {
			f, err = imageprocess.FormatFromExtension(fileURL)
			// 获取文件后缀
			if err != nil {
				// 将处理后的文件内容写入响应
				if _, err := io.Copy(w, file); err != nil {
					logx.Errorf("Copy file error: %s", err.Error())
					http.Error(w, "Failed to send file", http.StatusInternalServerError)
				}
				return
			}
		}
		//处理处理参数
		if len(parts) == 1 {
			//无需处理
			if _, err := io.Copy(w, file); err != nil {
				logx.Errorf("Copy file error: %s", err.Error())
				http.Error(w, "Failed to send file", http.StatusInternalServerError)
			}
			return
		}
		options, err := imageprocess.ParseOptions(strings.Split(strings.Split(parts[1], "?")[0], "&")[0])
		if err != nil {
			logx.Errorf("ParseOptions %s", err.Error())
			http.Error(w, fmt.Sprintf("ParseOptions %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if len(options) == 0 {
			//无需处理
			if _, err := io.Copy(w, file); err != nil {
				logx.Errorf("Copy file error: %s", err.Error())
				http.Error(w, "Failed to send file", http.StatusInternalServerError)
			}
			return
		}
		//处理图片
		err = processImg(file, w, f, options)
		if err != nil {
			logx.Errorf("processFile %s", err.Error())
			http.Error(w, fmt.Sprintf("processFile %s", err.Error()), http.StatusInternalServerError)
		}
	}
}

// 将url 转换成 md5
func httpCache(fileURL string, svcCtx *svc.ServiceContext) (*http.Response, io.ReadCloser, error) {
	key := strings.Replace(fileURL, "//", "/", -1)
	cache := svcCtx.Cache
	file, _, err := cache.GetStream(key)
	if err == nil {
		re := &http.Response{
			StatusCode: 200,
		}
		if svcCtx.Config.ImageCacheConf.Node == config.S3Node {
			headers, err := cache.HeadObject(key)
			if err == nil {
				re.Header = headers
				return re, file, nil
			}
		}
		typeData, _, err := cache.Get("header/" + key)
		if err == nil {
			headers := map[string][]string{}
			json.Unmarshal(typeData, &headers)
			// 构建http.Response
			re.Header = headers
			return re, file, nil
		}
		return re, file, nil
	}
	response, err := http.Get(fileURL)
	if err != nil {
		return nil, nil, err
	}
	// 读取响应体
	if response.ContentLength < 10*1024*1024 && response.StatusCode == 200 {
		// 缓存
		_, err = cache.PutStream(key, response.Body)
		if err != nil {
			return nil, nil, err
		}
		if svcCtx.Config.ImageCacheConf.Node != config.S3Node {
			headers := map[string][]string{}
			// 获取所有 header
			for key, value := range response.Header {
				headers[key] = value
			}
			data, _ := json.Marshal(headers)
			// 设置response信息
			cache.Put("header/"+key, data)
		}
		file, _, err = cache.GetStream(key)
	} else {
		file = response.Body
	}
	return response, file, err
}
