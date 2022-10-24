package jago

import (
	"log"
	"net/http"
)

type (
	Response struct {
		Writer    http.ResponseWriter
		Status    int
		Size      int64
		Committed bool
	}
)

func NewResponse(w http.ResponseWriter) (r *Response) {
	return &Response{Writer: w}
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) WriteHeader(code int) {
	if r.Committed {
		log.Println("response already committed")
		return
	}
	r.Status = code
	r.Writer.WriteHeader(r.Status)
	r.Committed = true
}

func (r *Response) SetHeader(key string, val string) {
	r.Writer.Header().Add(key, val)
}

func (r *Response) Write(b []byte) (n int, err error) {
	if !r.Committed {
		if r.Status == 0 {
			r.Status = http.StatusOK
		}
		r.WriteHeader(r.Status)
	}
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	return
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.Writer, cookie)
}

func (r *Response) Flush() {
	r.Writer.(http.Flusher).Flush()
}
