package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamesYYang/jago"
)

// /users/new
// /users/:id
// /users/:id/address/:address
// /users/:id/report/download
func main() {
	jago := jago.New()
	jago.Get("/users/new", TestHandler)
	jago.Get("/users/:id", TestHandler)
	jago.Get("/users/:id/address/:address", TestHandler)
	jago.Get("/users/:id/report/download", TestHandler)
	jago.Get("/users/*", TestHandler)
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: jago,
		// 请求监听地址
		Addr: ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func TestHandler(c jago.Context) error {
	result := "hello Jago\n"
	result += fmt.Sprintf("your path: %s \n", c.Request().URL)
	result += fmt.Sprintf("match path: %s \n", c.Path())
	result += fmt.Sprintf("params id: %s \n", c.Param("id"))
	result += fmt.Sprintf("params address: %s \n", c.Param("address"))
	return c.String(http.StatusOK, result)
}
