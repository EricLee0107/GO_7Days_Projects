package elee

import (
	"fmt"
	"net/http"
)

type Context interface{
	// Request returns `*http.Request`.
	Request() *http.Request

	// SetRequest sets `*http.Request`.
	SetRequest(r *http.Request)

	// SetResponse sets `*Response`.
	SetResponse(r *Response)

	// Response returns `*Response`.
	Response() *Response

	// Path returns the registered path for the handler.
	Path() string

	// SetPath sets the registered path for the handler.
	SetPath(p string)

	Method() string

	Header(key string)

	SetHeader(key string,value string)

	SetStatus(code int)

	FormValue (key string) string

	String(code int, format string,values ...interface{})


}

type HandlerFunc func(c Context)error


type context struct{
	request *http.Request // http请求
	response *Response // http请求响应，封装了http.ResponseWriter和状态码
	path string  // 请求路径
	method string // 请求的方法
}

func newContext(w http.ResponseWriter,req *http.Request) Context {
	res := &Response{}
	res.Writer = w
	cnt := &context{}
	cnt.request = req
	cnt.response = res
	cnt.path = req.URL.Path
	cnt.method = req.Method
	return cnt
}


func (cnt *context) String(code int, format string,values ...interface{}){
	cnt.SetHeader("Content-Type","text/plain")
	cnt.SetStatus(code)
	cnt.response.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

func (cnt *context) FormValue(key string) string{
	return cnt.request.FormValue(key)
}


func (cnt *context) SetStatus(code int){
	cnt.response.Writer.WriteHeader(code)
}

func (cnt *context)Header(key string){
	cnt.response.Writer.Header().Get(key)
}

func (cnt *context)SetHeader(key string,value string){
	cnt.response.Writer.Header().Set(key,value)
}


func (cnt *context) Method() string{
	return cnt.method
}


func(cnt *context) Request() *http.Request{
	return cnt.request
}

func (cnt *context) SetRequest(req *http.Request) {
	cnt.request = req
}

func (cnt *context) Response() *Response{
	return cnt.response
}

func (cnt *context) SetResponse(res *Response) {
	cnt.response = res
}

func (cnt *context) Path() string{
	return cnt.path
}

func (cnt *context) SetPath(path string) {
	cnt.path = path
}
