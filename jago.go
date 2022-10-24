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

func (j *Jago) Connect(path string, handler HandlerFunc) {
	j.Add(http.MethodConnect, path, handler)
}

func (j *Jago) Head(path string, handler HandlerFunc) {
	j.Add(http.MethodHead, path, handler)
}

func (j *Jago) Options(path string, handler HandlerFunc) {
	j.Add(http.MethodOptions, path, handler)
}

func (j *Jago) Patch(path string, handler HandlerFunc) {
	j.Add(http.MethodPatch, path, handler)
}

func (j *Jago) Trace(path string, handler HandlerFunc) {
	j.Add(http.MethodTrace, path, handler)
}

func (j *Jago) Get(path string, handler HandlerFunc) {
	j.Add(http.MethodGet, path, handler)
}

func (j *Jago) Post(path string, handler HandlerFunc) {
	j.Add(http.MethodPost, path, handler)
}

func (j *Jago) Put(path string, handler HandlerFunc) {
	j.Add(http.MethodPut, path, handler)
}

func (j *Jago) Delete(path string, handler HandlerFunc) {
	j.Add(http.MethodDelete, path, handler)
}

func (j *Jago) Group(prefix string) (g *Group) {
	g = &Group{prefix: prefix, j: j}
	return g
}

func (j *Jago) Add(method, path string, handler HandlerFunc) {
	j.router.add(method, path, handler)
}

func (j *Jago) findRoute(request *http.Request, c Context) {
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
