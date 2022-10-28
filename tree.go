package jago

import (
	"log"
	"strings"
)

type (
	Trie struct {
		root *TreeNode
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
		hasWildcard       bool
		handlers          map[string][]HandlerFunc
	}
)

func newTrie() *Trie {
	return &Trie{
		root: &TreeNode{
			parent:        nil,
			segChildren:   make(map[string]*TreeNode),
			paramChildren: make(map[string]*TreeNode),
		},
	}
}

func (t *Trie) add(method, pattern string, handlers ...HandlerFunc) {
	segments := getURIPaths(pattern)
	if len(segments) == 0 {
		return
	}
	node := parsePattern(t.root, segments)
	if node != nil {
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
	}

	if method == CONNECT || method == HttpMethodAny {
		node.handlers[CONNECT] = append(node.handlers[CONNECT], handlers...)
	} else if method == DELETE || method == HttpMethodAny {
		node.handlers[DELETE] = append(node.handlers[DELETE], handlers...)
	} else if method == GET || method == HttpMethodAny {
		node.handlers[GET] = append(node.handlers[GET], handlers...)
	} else if method == HEAD || method == HttpMethodAny {
		node.handlers[HEAD] = append(node.handlers[HEAD], handlers...)
	} else if method == OPTIONS || method == HttpMethodAny {
		node.handlers[OPTIONS] = append(node.handlers[OPTIONS], handlers...)
	} else if method == PATCH || method == HttpMethodAny {
		node.handlers[PATCH] = append(node.handlers[PATCH], handlers...)
	} else if method == POST || method == HttpMethodAny {
		node.handlers[POST] = append(node.handlers[POST], handlers...)
	} else if method == PUT || method == HttpMethodAny {
		node.handlers[PUT] = append(node.handlers[PUT], handlers...)
	} else if method == TRACE || method == HttpMethodAny {
		node.handlers[TRACE] = append(node.handlers[TRACE], handlers...)
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

func (t *Trie) printTree() {
	log.Println("root")
	prefix := ""
	prefix += "    "
	printNode(t.root, prefix)
}

func printNode(node *TreeNode, prefix string) {
	if len(node.segChildren) > 0 {
		for segment, n := range node.segChildren {
			log.Println(prefix + segment)
			printNode(n, prefix+"    ")
		}
	}

	if len(node.paramChildren) > 0 {
		for segment, n := range node.paramChildren {
			log.Println(prefix + segment)
			printNode(n, prefix+"    ")
		}
	}

	if node.wildcardChild != nil {
		log.Println(prefix + node.wildcardChild.segment)
	}
}
