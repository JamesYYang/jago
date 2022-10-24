package jago

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type (
	Context interface {
		Request() *http.Request
		Next() error

		Path() string
		Param(name string) string
		RealIP() string

		QueryParam(name string) string
		QueryParams() url.Values
		QueryString() string

		Cookie(name string) (*http.Cookie, error)
		Cookies() []*http.Cookie

		BindJson(i interface{}) error

		HTML(code int, html string) error
		HTMLBlob(code int, b []byte) error
		String(code int, s string) error

		JSON(code int, i interface{}) error
		JSONPretty(code int, i interface{}, indent string) error

		XML(code int, i interface{}) error
		XMLPretty(code int, i interface{}, indent string) error

		Blob(code int, contentType string, b []byte) error

		NoContent(code int) error
		Redirect(code int, url string) error
	}

	context struct {
		request  *http.Request
		response *Response
		j        *Jago
		path     string
		pnames   map[string]string
		query    url.Values
		handlers []HandlerFunc
		hIndex   int
	}
)

func (c *context) Request() *http.Request {
	return c.request
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

func (c *context) writeContentType(value string) {
	header := c.response.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *context) Path() string {
	return c.path
}

func (c *context) Param(name string) string {
	return c.pnames[name]
}

func (c *context) RealIP() string {
	// Fall back to legacy behavior
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			return strings.TrimSpace(ip[:i])
		}
		return ip
	}
	if ip := c.request.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.request.RemoteAddr)
	return ra
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

func (c *context) QueryString() string {
	return c.request.URL.RawQuery
}

func (c *context) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

func (c *context) BindJson(obj interface{}) error {
	err := json.NewDecoder(c.request.Body).Decode(obj)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset))
	} else if se, ok := err.(*json.SyntaxError); ok {
		return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
	}
	return err
}

func (c *context) HTML(code int, html string) (err error) {
	return c.HTMLBlob(code, []byte(html))
}

func (c *context) HTMLBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) json(code int, i interface{}, indent string) error {
	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	enc := json.NewEncoder(c.response)
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

func (c *context) JSONPretty(code int, i interface{}, indent string) (err error) {
	return c.json(code, i, indent)
}

func (c *context) xml(code int, i interface{}, indent string) (err error) {
	c.writeContentType(MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	enc := xml.NewEncoder(c.response)
	if indent != "" {
		enc.Indent("", indent)
	}
	if _, err = c.response.Write([]byte(xml.Header)); err != nil {
		return
	}
	return enc.Encode(i)
}

func (c *context) XML(code int, i interface{}) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.j.Debug || pretty {
		indent = defaultIndent
	}
	return c.xml(code, i, indent)
}

func (c *context) XMLPretty(code int, i interface{}, indent string) (err error) {
	return c.xml(code, i, indent)
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

func (c *context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

func (c *context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}
