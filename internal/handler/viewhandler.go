package handler

import (
	"fmt"
	"image/gif"
	"io"
	"net/http"
	"slimfiler/internal/svc"
	"slimfiler/internal/utils/fileutil"
	"strings"

	"github.com/AndsGo/imageprocess"
)

func ViewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只支持post
		if r.Method != http.MethodGet {
			http.Error(w, "only support post", http.StatusBadRequest)
			return
		}
		logx := svcCtx.Logger
		// 获取url path
		// 获取文件名称
		fileName := r.URL.Path
		// 打开文件
		file, _, err := svcCtx.Storage.GetStream(fileName)
		// file, err := os.Open(fmt.Sprintf("%s%s", svcCtx.Config.UploadConf.PublicStorePath, fileName))
		if err != nil {
			logx.Errorf("Open file error: %s", err.Error())
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()
		// 设置直接下载
		// 获取文件名
		extArr := strings.Split(fileName, "/")
		fileName = extArr[len(extArr)-1]
		fileutil.SetDownload(w, r, fileName)
		// 获取文件类型
		// 获取参数
		// 获取文件后缀
		f, err := imageprocess.FormatFromExtension(fileName)
		if err != nil {
			// 将处理后的文件内容写入响应
			if _, err := io.Copy(w, file); err != nil {
				logx.Errorf("Copy file error: %s", err.Error())
				http.Error(w, "Failed to send file", http.StatusInternalServerError)
			}
			return
		}
		//处理处理参数
		ossParams := r.URL.Query().Get("x-oss-process")
		if ossParams == "" {
			//无需处理
			if _, err := io.Copy(w, file); err != nil {
				logx.Errorf("Copy file error: %s", err.Error())
				http.Error(w, "Failed to send file", http.StatusInternalServerError)
			}
			return
		}
		options, err := imageprocess.ParseOptions(ossParams)
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

// 进行转换
func processImg(file io.Reader, w io.Writer, f imageprocess.Format, options []imageprocess.Option) error {
	if f == imageprocess.GIF {
		imgGif, err := gif.DecodeAll(file)
		if err != nil {
			return err
		}
		return imageprocess.ProcessGif(imgGif, w, options)
	} else {
		img, err := imageprocess.DecodeImage(file, f)
		if err != nil {
			return err
		}
		return imageprocess.Process(img, w, f, options)
	}
}
