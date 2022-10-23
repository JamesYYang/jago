package jago

import (
	"strings"
)

type (
	Router struct {
		routes map[int][]*node
	}

	node struct {
		path              string
		componentList     []string
		literalsToMatch   []string
		variablesNames    []string
		variableArgsCount int
		hasWildcard       bool
		method            string
		handler           HandlerFunc
	}
)

func newRouter() *Router {
	r := &Router{
		routes: map[int][]*node{},
	}

	for _, i := range []int{1, 2, 3, 4, 5, 6} {
		r.routes[i] = nil
	}

	return r
}

func (r *Router) add(method, path string, h HandlerFunc) {
	n := newNode(method, path, h)
	if n == nil {
		return
	}

	pathCount := len(n.componentList)
	if n.hasWildcard {
		r.routes[6] = append(r.routes[6], n)
	} else {
		if pathCount > 5 {
			pathCount = 5
		}
		r.routes[pathCount] = append(r.routes[pathCount], n)
	}
}

func (r *Router) find(uri string, method string, c Context) {
	ctx := c.(*context)

	pathParts := getURIPaths(uri)
	pathCount := len(pathParts)
	if pathCount > 5 {
		pathCount = 5
	}

	maxScore, n := r.findMaxScore(pathCount, pathParts, method)

	if maxScore == 0 {
		pathCount = 6
		maxScore, n = r.findMaxScore(pathCount, pathParts, method)
	}

	if maxScore > 0 {
		ctx.pnames = n.getPathParam(pathParts)
		ctx.handler = n.handler
		ctx.path = n.path
	}
}

func (r *Router) findMaxScore(pathCount int, pathParts []string, method string) (int, *node) {
	maxScore := 0
	var n *node

	for _, route := range r.routes[pathCount] {
		score := route.matchScore(pathParts, method)
		if score > maxScore {
			maxScore = score
			n = route
		}
		if maxScore == 1001 {
			break
		}
	}
	return maxScore, n
}

func newNode(method string, uri string, h HandlerFunc) *node {
	cl := getURIPaths(uri)
	if len(cl) == 0 {
		return nil
	}
	n := &node{
		path:              uri,
		componentList:     cl,
		variableArgsCount: 0,
		hasWildcard:       false,
		handler:           h,
		method:            method,
	}
	componentLength := len(n.componentList)
	if n.componentList[componentLength-1] == "*" {
		n.componentList = n.componentList[:componentLength-1]
		n.hasWildcard = true
	}
	n.literalsToMatch = make([]string, componentLength)
	n.variablesNames = make([]string, componentLength)
	for i, component := range n.componentList {
		if strings.Index(component, ":") == 0 {
			n.variablesNames[i] = component[1:]
			n.variableArgsCount++
		} else {
			n.literalsToMatch[i] = strings.ToLower(component)
		}
	}
	return n
}

func (n *node) matchScore(pathParts []string, method string) int {
	if !n.hasWildcard && len(pathParts) != len(n.componentList) {
		return -1
	}
	for i := range n.componentList {
		if n.literalsToMatch[i] != "" {
			if i >= len(pathParts) || strings.ToLower(pathParts[i]) != n.literalsToMatch[i] {
				return -1
			}
		}
	}

	if n.method != HttpMethodAny && n.method != method {
		return -1
	}

	score := 1
	score += max(10-n.variableArgsCount, 1) * 100
	if n.hasWildcard {
		score += len(n.componentList)
	}
	return score
}

func (n *node) getPathParam(pathParts []string) map[string]string {
	pathParam := make(map[string]string)
	for i, pname := range n.variablesNames {
		if pname != "" {
			pathParam[pname] = pathParts[i]
		}
	}
	return pathParam
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
