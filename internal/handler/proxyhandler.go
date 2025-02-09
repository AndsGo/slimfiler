package handler

import (
	"fileserver/internal/svc"
	"fileserver/internal/utils/fileutil"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/AndsGo/imageprocess"
	"github.com/google/uuid"
)

func ProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只支持post
		if r.Method != http.MethodGet {
			http.Error(w, "only support post", http.StatusBadRequest)
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
		response, err := http.Get(fileURL)
		if err != nil {
			fmt.Println("Error while making request:", err)
			return
		}
		// response to Reader
		file := io.Reader(response.Body)
		defer response.Body.Close()
		// 获取 response 文件类型
		fileType := strings.Split(response.Header.Get("Content-Type"), "/")[0]
		if fileType != "image" && fileType != "video" && fileType != "audio" {
			fileType = "other"
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
