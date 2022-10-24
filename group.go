package jago

import "net/http"

type (
	// Group is a set of sub-routes for a specified route.
	Group struct {
		j      *Jago
		prefix string
	}
)

func (g *Group) Connect(path string, handler HandlerFunc) {
	g.Add(http.MethodConnect, path, handler)
}

func (g *Group) Head(path string, handler HandlerFunc) {
	g.Add(http.MethodHead, path, handler)
}

func (g *Group) Options(path string, handler HandlerFunc) {
	g.Add(http.MethodOptions, path, handler)
}

func (g *Group) Patch(path string, handler HandlerFunc) {
	g.Add(http.MethodPatch, path, handler)
}

func (g *Group) Trace(path string, handler HandlerFunc) {
	g.Add(http.MethodTrace, path, handler)
}

func (g *Group) Get(path string, handler HandlerFunc) {
	g.Add(http.MethodGet, path, handler)
}

func (g *Group) Post(path string, handler HandlerFunc) {
	g.Add(http.MethodPost, path, handler)
}

func (g *Group) Put(path string, handler HandlerFunc) {
	g.Add(http.MethodPut, path, handler)
}

func (g *Group) Delete(path string, handler HandlerFunc) {
	g.Add(http.MethodDelete, path, handler)
}

func (g *Group) Add(method, path string, handler HandlerFunc) {
	g.j.Add(method, g.prefix+path, handler)
}
