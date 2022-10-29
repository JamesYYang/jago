package jago

import (
	"log"
	"strings"
)

type (
	Trie struct {
		root           *TreeNode
		staticChildren map[string]*TreeNode
	}

	TreeNode struct {
		parent            *TreeNode
		segment           string
		wildcardChild     *TreeNode
		segChildren       map[string]*TreeNode
		paramChildren     map[string]*TreeNode
		leaf              bool
		path              string
		componentList     []string
		literalsToMatch   []string
		variablesNames    []string
		variableArgsCount int
		score             int
		hasWildcard       bool
		handlers          map[string][]HandlerFunc
	}

	Mached struct {
		results []*TreeNode
	}
)

func newTrie() *Trie {
	return &Trie{
		root: &TreeNode{
			parent:        nil,
			segChildren:   make(map[string]*TreeNode),
			paramChildren: make(map[string]*TreeNode),
		},
		staticChildren: make(map[string]*TreeNode),
	}
}

func (t *Trie) add(method, pattern string, handlers ...HandlerFunc) {
	pattern = strings.ToLower(pattern)
	if pattern == "/" {
		pattern = "/*"
	}
	segments := getURIPaths(pattern)
	if len(segments) == 0 {
		return
	}

	var node *TreeNode
	if !strings.Contains(pattern, ":") && !strings.Contains(pattern, "*") {
		node = &TreeNode{
			hasWildcard: false,
			segment:     pattern,
			handlers:    make(map[string][]HandlerFunc),
		}
		t.staticChildren[pattern] = node
	} else {
		node = parsePattern(t.root, segments)
	}
	if node != nil {
		initLeafNode(node, method, pattern, segments, handlers...)
	}
}

func initLeafNode(node *TreeNode, method, pattern string, segments []string, handlers ...HandlerFunc) {
	node.path = pattern
	node.componentList = segments
	componentLength := len(node.componentList)
	if node.hasWildcard {
		node.componentList = node.componentList[:componentLength-1]
	}
	node.literalsToMatch = make([]string, componentLength)
	node.variablesNames = make([]string, componentLength)
	for i, component := range node.componentList {
		if strings.Index(component, ":") == 0 {
			node.variablesNames[i] = component[1:]
			node.variableArgsCount++
		} else {
			node.literalsToMatch[i] = strings.ToLower(component)
		}
	}

	if method == CONNECT || method == HttpMethodAny {
		node.handlers[CONNECT] = append(node.handlers[CONNECT], handlers...)
	}
	if method == DELETE || method == HttpMethodAny {
		node.handlers[DELETE] = append(node.handlers[DELETE], handlers...)
	}
	if method == GET || method == HttpMethodAny {
		node.handlers[GET] = append(node.handlers[GET], handlers...)
	}
	if method == HEAD || method == HttpMethodAny {
		node.handlers[HEAD] = append(node.handlers[HEAD], handlers...)
	}
	if method == OPTIONS || method == HttpMethodAny {
		node.handlers[OPTIONS] = append(node.handlers[OPTIONS], handlers...)
	}
	if method == PATCH || method == HttpMethodAny {
		node.handlers[PATCH] = append(node.handlers[PATCH], handlers...)
	}
	if method == POST || method == HttpMethodAny {
		node.handlers[POST] = append(node.handlers[POST], handlers...)
	}
	if method == PUT || method == HttpMethodAny {
		node.handlers[PUT] = append(node.handlers[PUT], handlers...)
	}
	if method == TRACE || method == HttpMethodAny {
		node.handlers[TRACE] = append(node.handlers[TRACE], handlers...)
	}

	node.score = 1
	baseScore := 100
	if node.hasWildcard {
		baseScore = 10
	}
	node.score += max(10-node.variableArgsCount, 1) * baseScore
	if node.hasWildcard {
		node.score += len(node.componentList)
	}
}

func parsePattern(parent *TreeNode, segments []string) *TreeNode {
	segment := segments[0]
	segments = segments[1:]
	child := parseSegment(parent, segment)
	if len(segments) == 0 {
		child.leaf = true
		return child
	}
	return parsePattern(child, segments)
}

func parseSegment(parent *TreeNode, segment string) *TreeNode {
	if node, ok := parent.segChildren[segment]; ok {
		return node
	}
	if strings.HasPrefix(segment, ":") {
		if node, ok := parent.paramChildren[segment]; ok {
			return node
		}
	}
	node := &TreeNode{
		segment:       segment,
		parent:        parent,
		segChildren:   make(map[string]*TreeNode),
		paramChildren: make(map[string]*TreeNode),
		handlers:      make(map[string][]HandlerFunc),
	}
	if segment == "*" {
		parent.wildcardChild = node
		node.hasWildcard = true
	} else if strings.HasPrefix(segment, ":") {
		parent.paramChildren[segment] = node
	} else {
		parent.segChildren[segment] = node
	}
	return node
}

func (t *Trie) find(uri, method string) (maxScore int, node *TreeNode) {
	uri = strings.ToLower(uri)
	if uri == "/" {
		uri = "/*"
	}

	if n, ok := t.staticChildren[uri]; ok {
		if _, ok := n.handlers[method]; ok {
			maxScore = n.score
			node = n
			return
		}
	}

	pathParts := getURIPaths(uri)
	matched := &Mached{
		results: make([]*TreeNode, 0),
	}
	matchNode(t.root, method, pathParts, matched)

	for _, n := range matched.results {
		// log.Printf("uri: %s, matched node: %s", uri, n.path)
		if n.score > maxScore {
			if _, ok := n.handlers[method]; ok {
				maxScore = n.score
				node = n
			}
		}
	}
	return
}

func matchNode(parent *TreeNode, method string, pathParts []string, m *Mached) {
	segment := pathParts[0]
	segments := pathParts[1:]

	if len(segments) == 0 {

		if n, ok := parent.segChildren[segment]; ok {
			if n.leaf {
				m.results = append(m.results, n)
			}
		}

		for _, n := range parent.paramChildren {
			if n.leaf {
				m.results = append(m.results, n)
			}
		}

		if parent.wildcardChild != nil {
			m.results = append(m.results, parent.wildcardChild)
		}

	} else {

		if n, ok := parent.segChildren[segment]; ok {
			matchNode(n, method, segments, m)
		}

		for _, n := range parent.paramChildren {
			matchNode(n, method, segments, m)
		}

		if parent.wildcardChild != nil {
			m.results = append(m.results, parent.wildcardChild)
		}

	}
}

func (t *Trie) printTree() {
	log.Println("root")
	prefix := ""
	prefix += "    "
	for s, n := range t.staticChildren {
		log.Printf("%s%s [%d] -- %v", prefix, s, n.score, n.leaf)
	}
	printNode(t.root, prefix)
}

func printNode(node *TreeNode, prefix string) {
	if len(node.segChildren) > 0 {
		for segment, n := range node.segChildren {
			log.Printf("%s%s [%d] -- %v", prefix, segment, n.score, n.leaf)
			printNode(n, prefix+"    ")
		}
	}

	if len(node.paramChildren) > 0 {
		for segment, n := range node.paramChildren {
			log.Printf("%s%s [%d] -- %v", prefix, segment, n.score, n.leaf)
			printNode(n, prefix+"    ")
		}
	}

	if node.wildcardChild != nil {
		log.Printf("%s%s [%d] -- %v", prefix, node.wildcardChild.segment, node.wildcardChild.score, node.wildcardChild.leaf)
	}
}

func (n *TreeNode) getPathParam(pathParts []string) map[string]string {
	pathParam := make(map[string]string)
	for i, pname := range n.variablesNames {
		if pname != "" {
			pathParam[pname] = pathParts[i]
		}
	}
	return pathParam
}
