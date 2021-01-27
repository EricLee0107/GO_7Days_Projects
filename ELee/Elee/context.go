package Elee

import "net/http"

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

	// Handler returns the matched handler by router.
	Handler() HandlerFunc

	// SetHandler sets the matched handler by router.
	SetHandler(h HandlerFunc)

}

type HandlerFunc func(c *Context)error


type context struct{
	request *http.Request // http请求
	response *Response // http请求响应，封装了http.ResponseWriter和状态码
	path string  // 请求路径
	method string
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
func (cnt *context) Handler() HandlerFunc{
	return cnt.handler
}
func (cnt *context) SetHandler(h HandlerFunc){
	cnt.handler = h
}
