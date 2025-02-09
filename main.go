package main

import (
	"fileserver/internal/config"
	"fileserver/internal/handler"
	"fileserver/internal/svc"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// 文件夹，有need change it
var configFile = flag.String("f", "./etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	conf := loadConfig(*configFile)
	// 日志设置
	svcContxet := svc.NewServiceContext(conf)
	svcContxet.Logger.Info("Starting server")
	http.HandleFunc("/upload", handler.UploadHandler(svcContxet))
	http.HandleFunc("/proxy/", handler.ProxyHandler(svcContxet))
	http.HandleFunc("/", handler.ViewHandler(svcContxet))

	fmt.Println(fmt.Printf("Starting server on :%d", conf.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}

func loadConfig(configFile string) *config.Config {
	conf := config.Config{}
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config.yaml: %v", err)
	}
	fmt.Println(string(data))
	if err = yaml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("failed to parse config.yaml: %v", err)
	}
	return &conf
}
