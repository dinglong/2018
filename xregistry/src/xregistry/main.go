package main

import (
	"log"
	"net/http"

	_ "xregistry/frontend"
	"xregistry/handler"
	"xregistry/util"
)

func main() {
	http.DefaultServeMux.HandleFunc("/file", handler.LoadFileHandler)

	// 启动浏览器
	go util.OpenBrowser("http://127.0.0.1:8080")

	// 启动服务
	log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}
