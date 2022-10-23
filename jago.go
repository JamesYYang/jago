package jago

import (
	"fmt"
	"log"
	"net/http"
)

type (
	Jago struct {
		router *Router
	}

	HandlerFunc func(c Context) error
)

func New() *Jago {
	log.Printf(banner, Version)
	return &Jago{
		router: newRouter(),
	}
}

func (j *Jago) NewContext(r *http.Request, w http.ResponseWriter) Context {
	return &context{
		request:        r,
		responseWriter: w,
		j:              j,
	}
}

func (j *Jago) Get(url string, handler HandlerFunc) {
	j.router.add(http.MethodGet, url, handler)
}

func (j *Jago) Post(url string, handler HandlerFunc) {
	j.router.add(http.MethodPost, url, handler)
}

func (j *Jago) Put(url string, handler HandlerFunc) {
	j.router.add(http.MethodPut, url, handler)
}

func (j *Jago) Delete(url string, handler HandlerFunc) {
	j.router.add(http.MethodDelete, url, handler)
}

func (j *Jago) findRoute(request *http.Request, c Context) {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method

	j.router.find(uri, method, c)
}

func (j *Jago) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("Jago serveHTTP")
	ctx := j.NewContext(request, response)

	j.findRoute(request, ctx)

	h := ctx.Handler()
	if h != nil {
		h(ctx)
	} else {
		ctx.String(http.StatusNotFound, fmt.Sprintf("your path: %s not found", request.URL))
	}
}
