package handler

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"slimfiler/internal/svc"
	md5util "slimfiler/internal/utils/md5"
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
			http.Error(w, "reject the file which does not have suffix", http.StatusBadRequest)
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
		if directory == "" {
			directory = path.Join(svcCtx.Config.Name, fileType, timeString)
		}
		directory = path.Join("/", directory)

		var md5 string = ""
		if r.MultipartForm.Value["md5"] != nil {
			md5 = r.MultipartForm.Value["md5"][0]
		}
		// 判断是否已经上传过
		shouldReturn := returnFormMd5(md5, handler.Filename, svcCtx, w)
		if shouldReturn {
			return
		}
		context, err := io.ReadAll(file)
		if err != nil {
			logx.Errorf("fail to read file content: %v", err)
			http.Error(w, "fail to read file content", http.StatusInternalServerError)
			return
		}
		if md5 == "" {
			md5 = md5util.GetMD5(context)
		}
		// 判断是否已经上传过
		shouldReturn = returnFormMd5(md5, handler.Filename, svcCtx, w)
		if shouldReturn {
			return
		}
		storeFileName = path.Join(directory, storeFileName)
		etag, err := svcCtx.Storage.Put(storeFileName, context)
		if err != nil {
			logx.Errorf("fail to put file to storage: %v", err)
			http.Error(w, "fail to put file to storage", http.StatusInternalServerError)
			return
		}
		// 保存到数据库
		svcCtx.Db.Set(md5, storeFileName)
		// 返回json
		returnJson(w, etag, handler.Filename, svcCtx, storeFileName)
	}
}

// returnFormMd5 根据MD5值返回表单信息。
// 该函数尝试从数据库中获取与给定MD5值关联的文件名，并调用returnJson函数返回相关信息。
// 如果找不到文件名，则返回false；如果成功找到并处理，则返回true。
//
// 参数:
//
//	md5 - 表单的MD5值，用于查询数据库中的文件名。
//	svcCtx - 服务上下文，包含数据库连接等信息。
//	w - HTTP响应写入器，用于向客户端发送响应。
//
// 返回值:
//
//	成功找到并处理文件名时返回true，否则返回false。
func returnFormMd5(md5 string, fileName string, svcCtx *svc.ServiceContext, w http.ResponseWriter) bool {
	// 检查MD5值是否为空，如果为空则尝试从数据库中获取文件名。
	if md5 != "" {
		fullfileName := ""
		// 从数据库中获取与MD5值关联的文件名。
		svcCtx.Db.Get(md5, &fullfileName)
		// 如果找到了文件名，则调用returnJson函数返回相关信息，并返回true。
		if fullfileName != "" {
			returnJson(w, md5, fileName, svcCtx, fullfileName)
			return true
		}
	}
	// 如果没有找到文件名，返回false。
	return false
}

// returnJson 返回JSON响应。
// 该函数负责向客户端返回包含文件信息的JSON对象，包括文件名和文件URL。
// 参数:
//
//	w: http.ResponseWriter，用于写入HTTP响应的接口。
//	etag: 字符串，表示文件的唯一标识符，用于HTTP缓存控制。
//	filename: 字符串，表示原始文件名。
//	svcCtx: *svc.ServiceContext，服务上下文，包含配置信息等。
//	storeFileName: 字符串，表示存储在服务器上的文件名。
func returnJson(w http.ResponseWriter, etag string, filename string, svcCtx *svc.ServiceContext, storeFileName string) {
	// 设置响应头，指定内容类型为JSON。
	w.Header().Set("Content-Type", "application/json")
	// 允许跨域请求。
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 设置ETag，用于标识文件。
	w.Header().Set("ETag", etag)
	// 发送HTTP状态码200，表示请求成功。
	w.WriteHeader(http.StatusOK)
	// 构造JSON响应体，包含文件名和文件URL，并写入响应中。
	w.Write([]byte(fmt.Sprintf(`{"code":10000,"data":{"name":"%s","url":"%s"},"msg":""}`, filename, svcCtx.Config.UploadConf.ServerURL+storeFileName)))
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
