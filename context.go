package jago

import "net/http"

type (
	Context interface {
		Request() *http.Request
		Handler() HandlerFunc
		Path() string
		Param(name string) string
		String(code int, s string) error
	}

	context struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		j              *Jago
		path           string
		pnames         map[string]string
		handler        HandlerFunc
	}
)

func (c *context) writeContentType(value string) {
	header := c.responseWriter.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *context) Param(name string) string {
	return c.pnames[name]
}

func (c *context) Path() string {
	return c.path
}

func (c *context) Handler() HandlerFunc {
	return c.handler
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.responseWriter.WriteHeader(code)
	_, err = c.responseWriter.Write(b)
	return
}
