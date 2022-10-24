package jago

import "net/http"

type (
	// Group is a set of sub-routes for a specified route.
	Group struct {
		j           *Jago
		prefix      string
		middlewares []HandlerFunc
	}
)

func (g *Group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Connect(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodConnect, path, handlers...)
}

func (g *Group) Head(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodHead, path, handlers...)
}

func (g *Group) Options(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodOptions, path, handlers...)
}

func (g *Group) Patch(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodPatch, path, handlers...)
}

func (g *Group) Trace(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodTrace, path, handlers...)
}

func (g *Group) Get(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodGet, path, handlers...)
}

func (g *Group) Post(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodPost, path, handlers...)
}

func (g *Group) Put(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodPut, path, handlers...)
}

func (g *Group) Delete(path string, handlers ...HandlerFunc) {
	g.Add(http.MethodDelete, path, handlers...)
}

func (g *Group) Any(path string, handlers ...HandlerFunc) {
	g.Add(HttpMethodAny, path, handlers...)
}

func (g *Group) Add(method, path string, handlers ...HandlerFunc) {
	allHandlers := append(g.middlewares, handlers...)
	g.j.Add(method, g.prefix+path, allHandlers...)
}
