package jago

import (
	"fmt"
	"html"
	"net/http"
)

type Jago struct {
}

func New() *Jago {
	return &Jago{}
}

func (j *Jago) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, %q", html.EscapeString(request.URL.Path))
}
