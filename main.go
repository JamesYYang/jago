package main

import (
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: NewJago(),
		// 请求监听地址
		Addr: ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
