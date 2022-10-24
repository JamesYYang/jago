package jago

import (
	"fmt"
	"log"
	"net/http"
)

type (
	Jago struct {
		router           *Router
		middlewares      []HandlerFunc
		HTTPErrorHandler HTTPErrorHandler
		Debug            bool
	}

	HTTPError struct {
		Code    int         `json:"-"`
		Message interface{} `json:"message"`
	}

	HandlerFunc      func(c Context) error
	HTTPErrorHandler func(error, Context)
)

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

func New() *Jago {
	log.Printf(banner, Version)
	j := &Jago{
		router: newRouter(),
	}
	j.HTTPErrorHandler = j.DefaultHTTPErrorHandler

	return j
}

func (j *Jago) NewContext(r *http.Request, w http.ResponseWriter) Context {
	return &context{
		request:  r,
		response: NewResponse(w),
		j:        j,
		hIndex:   -1,
	}
}

func (j *Jago) Use(middlewares ...HandlerFunc) {
	j.middlewares = append(j.middlewares, middlewares...)
}

func (j *Jago) Connect(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodConnect, path, handlers...)
}

func (j *Jago) Head(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodHead, path, handlers...)
}

func (j *Jago) Options(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodOptions, path, handlers...)
}

func (j *Jago) Patch(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodPatch, path, handlers...)
}

func (j *Jago) Trace(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodTrace, path, handlers...)
}

func (j *Jago) Get(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodGet, path, handlers...)
}

func (j *Jago) Post(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodPost, path, handlers...)
}

func (j *Jago) Put(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodPut, path, handlers...)
}

func (j *Jago) Delete(path string, handlers ...HandlerFunc) {
	j.Add(http.MethodDelete, path, handlers...)
}

func (j *Jago) Any(path string, handlers ...HandlerFunc) {
	j.Add(HttpMethodAny, path, handlers...)
}

func (j *Jago) Group(prefix string, handlers ...HandlerFunc) (g *Group) {
	g = &Group{prefix: prefix, j: j}
	g.Use(handlers...)
	return g
}

func (j *Jago) Add(method, path string, handlers ...HandlerFunc) {
	allHandlers := append(j.middlewares, handlers...)
	j.router.add(method, path, allHandlers...)
}

func (j *Jago) findRoute(request *http.Request, c Context) {
	uri := request.URL.Path
	method := request.Method

	j.router.find(uri, method, c)
}

func (j *Jago) DefaultHTTPErrorHandler(err error, c Context) {
	he, ok := err.(*HTTPError)
	if !ok {
		he = &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	message := he.Message

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		log.Println(err)
	}
}

func (j *Jago) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("Jago serveHTTP")
	ctx := j.NewContext(request, response)

	j.findRoute(request, ctx)
	if err := ctx.Next(); err != nil {
		j.HTTPErrorHandler(err, ctx)
	}
}
