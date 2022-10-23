package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamesYYang/jago"
)

func main() {
	jago := jago.New()
	jago.Get("foo", FooControllerHandler)
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: jago,
		// 请求监听地址
		Addr: ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func FooControllerHandler(c jago.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("hello, your path: %s", c.Request().URL))
}
