package jago

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type (
	Context interface {
		Request() *http.Request
		Next() error
		Path() string
		Param(name string) string
		QueryParam(name string) string
		QueryParams() url.Values
		String(code int, s string) error
		JSON(code int, i interface{}) error
		NoContent(code int) error
	}

	context struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		j              *Jago
		path           string
		pnames         map[string]string
		query          url.Values
		handlers       []HandlerFunc
		hIndex         int
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

func (c *context) Next() error {
	c.hIndex++
	if c.hIndex < len(c.handlers) {
		if err := c.handlers[c.hIndex](c); err != nil {
			return err
		}
	}
	return nil
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) QueryParam(name string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query.Get(name)
}

func (c *context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query
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

func (c *context) NoContent(code int) error {
	c.responseWriter.WriteHeader(code)
	return nil
}

func (c *context) json(code int, i interface{}, indent string) error {
	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	c.responseWriter.WriteHeader(code)
	enc := json.NewEncoder(c.responseWriter)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

func (c *context) JSON(code int, i interface{}) error {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.j.Debug || pretty {
		indent = defaultIndent
	}
	return c.json(code, i, indent)
}
