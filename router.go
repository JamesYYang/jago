package jago

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"

	tlcache "github.com/JamesYYang/go-ttl-lru"
)

type (
	Router struct {
		routes map[string]map[int][]*node
		cache  *tlcache.Cache
	}

	node struct {
		path              string
		componentList     []string
		literalsToMatch   []string
		variablesNames    []string
		variableArgsCount int
		hasWildcard       bool
		method            string
		handlers          []HandlerFunc
	}
)

func newRouter() *Router {
	r := &Router{
		routes: make(map[string]map[int][]*node),
	}
	r.cache = tlcache.NewLRUCache(100)

	// for _, i := range []int{1, 2, 3, 4, 5, 6} {
	// 	r.routes[i] = nil
	// }

	return r
}

func (r *Router) add(method, path string, handlers ...HandlerFunc) {
	n := newNode(method, path, handlers...)
	if n == nil {
		return
	}

	pathCount := len(n.componentList)
	pathKey := strings.ToLower(n.literalsToMatch[0])
	if pathKey == "" {
		pathKey = "#dynamic#"
	}
	pathTree, ok := r.routes[pathKey]
	if !ok {
		newTree := make(map[int][]*node)
		for _, i := range []int{1, 2, 3, 4, 5, 6} {
			newTree[i] = nil
		}
		r.routes[pathKey] = newTree
		pathTree = r.routes[pathKey]
	}

	if n.hasWildcard {
		pathTree[6] = append(pathTree[6], n)
	} else {
		if pathCount > 5 {
			pathCount = 5
		}
		pathTree[pathCount] = append(pathTree[pathCount], n)
	}
}

func (r *Router) findRoute(pathParts []string, method string) (maxScore int, n *node) {

	if n == nil {
		pathCount := len(pathParts)
		if pathCount > 5 {
			pathCount = 5
		}

		pathKey := strings.ToLower(pathParts[0])
		maxScore, n = r.findMaxScore(pathKey, pathCount, pathParts, method)

		if maxScore == 0 {
			maxScore, n = r.findMaxScore("#dynamic#", pathCount, pathParts, method)
		}

		if maxScore == 0 {
			pathCount = 6
			maxScore, n = r.findMaxScore(pathKey, pathCount, pathParts, method)
		}

		if maxScore == 0 {
			pathCount = 6
			maxScore, n = r.findMaxScore("#dynamic#", pathCount, pathParts, method)
		}

	}
	return
}

func (r *Router) find(uri string, method string, c Context) {
	ctx := c.(*context)

	// cacheKey := hash(uri)

	// maxScore, n = r.findNodeFromCache(cacheKey)

	pathParts := getURIPaths(uri)
	maxScore, n := r.findRoute(pathParts, method)

	// if maxScore > 0 {
	// 	r.setNode2Cache(cacheKey, n)
	// }

	if maxScore > 0 {
		ctx.pnames = n.getPathParam(pathParts)
		ctx.handlers = n.handlers
		ctx.path = n.path
	} else {
		ctx.handlers = append(ctx.handlers, NotFoundHandler)
	}
}

func (r *Router) findNodeFromCache(cacheKey string) (int, *node) {
	if v, ok := r.cache.Get(cacheKey); ok {
		return 1, v.(*node)
	} else {
		return 0, nil
	}
}

func (r *Router) setNode2Cache(cacheKey string, n *node) {
	r.cache.Add(cacheKey, n)
}

func (r *Router) findMaxScore(pathKey string, pathCount int, pathParts []string, method string) (int, *node) {
	maxScore := 0
	var n *node

	pathTree, ok := r.routes[pathKey]
	if !ok {
		return maxScore, n
	}

	for _, route := range pathTree[pathCount] {
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

func newNode(method string, uri string, handlers ...HandlerFunc) *node {
	cl := getURIPaths(uri)
	if len(cl) == 0 {
		return nil
	}
	n := &node{
		path:              uri,
		componentList:     cl,
		variableArgsCount: 0,
		hasWildcard:       false,
		handlers:          handlers,
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

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}
