package jago

import (
	"strings"
)

type (
	Router struct {
		routes *Trie
	}
)

func newRouter() *Router {
	r := &Router{
		routes: newTrie(),
	}

	return r
}

func (r *Router) add(method, path string, handlers ...HandlerFunc) {
	r.routes.add(method, path, handlers...)
}

func (r *Router) PrintTree() {
	r.routes.printTree()
}

func (r *Router) find(uri string, method string, c Context) {
	ctx := c.(*context)
	pathParts := getURIPaths(uri)
	maxScore, n := r.routes.find(uri, method)

	if maxScore > 0 {
		ctx.pnames = n.getPathParam(pathParts)
		ctx.handlers = n.handlers[method]
		ctx.path = n.path
	} else {
		ctx.handlers = append(ctx.handlers, NotFoundHandler)
	}
}

func getURIPaths(url string) []string {
	paths := strings.Split(url, "/")
	return filter(paths, func(v string) bool {
		return v != ""
	})
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
