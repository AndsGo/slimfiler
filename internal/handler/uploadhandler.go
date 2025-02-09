package handler

import (
	"fileserver/internal/svc"
	"fileserver/internal/utils/fileutil"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只支持post
		if r.Method != http.MethodPost {
			http.Error(w, "only support post", http.StatusBadRequest)
			return
		}
		logx := svcCtx.Logger
		err := r.ParseMultipartForm(svcCtx.Config.UploadConf.MaxVideoSize)
		if err != nil {
			logx.Error("fail to parse the multipart form")
			http.Error(w, "fail to parse the multipart form", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			logx.Error("the value of file cannot be found")
			http.Error(w, "the value of file cannot be found", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// judge if the suffix is legal
		// 校验后缀是否合法
		dotIndex := strings.LastIndex(handler.Filename, ".")
		// if there is no suffix, reject it
		// 拒绝无后缀文件
		if dotIndex == -1 {
			logx.Error("reject the file which does not have suffix")
			http.Error(w, "the value of file cannot be found", http.StatusBadRequest)
			return
		}

		_, fileSuffix := handler.Filename[:dotIndex], handler.Filename[dotIndex+1:]
		fileUUID := uuid.New()
		storeFileName := fileUUID.String() + "." + fileSuffix
		timeString := time.Now().Format("2006-01-02")
		// judge if the file size is over max size
		// 判断文件大小是否超过设定值
		fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[0]
		if fileType != "image" && fileType != "video" && fileType != "audio" {
			fileType = "other"
		}
		err = CheckOverSize(svcCtx, fileType, handler.Size)
		if err != nil {
			logx.Error("the file is over size")
			http.Error(w, "the file is over size", http.StatusBadRequest)
		}
		//
		var directory string
		if r.MultipartForm.Value["directory"] != nil {
			directory = r.MultipartForm.Value["directory"][0]
		}
		if directory != "" {
			directory = path.Join(directory, timeString)
		} else {
			directory = path.Join(svcCtx.Config.Name, fileType, timeString)
		}
		// generate path
		// 生成路径
		fullPath := path.Join(svcCtx.Config.UploadConf.PublicStorePath, directory)
		if !fileutil.IsExist(fullPath) {
			err = fileutil.CreateDir(fullPath + "/")
			if err != nil {
				logx.Error("failed to create directory for storing public files")
				http.Error(w, "failed to create directory for storing public files", http.StatusInternalServerError)
				return
			}
		}

		// default is public
		// 默认是公开的
		targetFile, err := os.Create(path.Join(fullPath, storeFileName))
		if err != nil {
			logx.Errorf("fail to create file,%v", err)
			http.Error(w, "fail to create file", http.StatusInternalServerError)
			return
		}
		defer targetFile.Close()
		_, err = io.Copy(targetFile, file)
		if err != nil {
			logx.Errorf("fail to copy file content: %v", err)
			http.Error(w, "fail to copy file content", http.StatusInternalServerError)
			return
		}
		// var md5 string
		// if r.MultipartForm.Value["md5"] != nil {
		// 	md5 = r.MultipartForm.Value["md5"][0]
		// } else {
		// 	md5 = ""
		// }
		// 返回json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"code":10000,"data":{"name":"%s","url":"%s"},"msg":""}`, handler.Filename, svcCtx.Config.UploadConf.ServerURL+"/"+directory+"/"+storeFileName)))
	}
}

func CheckOverSize(svCtx *svc.ServiceContext, fileType string, size int64) error {
	if fileType == "image" && size > svCtx.Config.UploadConf.MaxImageSize {
		return fmt.Errorf("file.overSizeError")
	} else if fileType == "video" && size > svCtx.Config.UploadConf.MaxVideoSize {
		return fmt.Errorf("file.overSizeError")
	} else if fileType == "audio" && size > svCtx.Config.UploadConf.MaxAudioSize {
		return fmt.Errorf("file.overSizeError")
	} else if fileType != "image" && fileType != "video" && fileType != "audio" &&
		size > svCtx.Config.UploadConf.MaxOtherSize {
		return fmt.Errorf("file.overSizeError")
	}
	return nil
}
