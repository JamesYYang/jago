package jago

import (
	"log"
	"net/http"
)

type (
	Jago struct {
		router map[string]HandlerFunc
	}

	HandlerFunc func(c Context) error
)

func New() *Jago {
	log.Printf(banner, Version)
	return &Jago{
		router: map[string]HandlerFunc{},
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
	j.router[url] = handler
}

func (j *Jago) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("Jago serveHTTP")
	ctx := j.NewContext(request, response)

	router := j.router["foo"]
	if router == nil {
		return
	}
	log.Println("Jago router")

	router(ctx)
}
