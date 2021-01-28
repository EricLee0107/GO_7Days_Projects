package elee

import (
	"net/http"
)



type Elee struct{
	// 静态路由
	routers map[string] *methodHandler
	tree *node
	maxParam *int
}


func New() *Elee{
	elee := &Elee{
		routers: make(map[string] *methodHandler),
	}
	return elee
}

func (elee *Elee) addRouter(method, path string, handler HandlerFunc){
	meth,ok := elee.routers[path]
	if !ok{
		meth = &methodHandler{}
		elee.routers[path] = meth
	}
	meth.addHandler(method,handler)
}


func (elee *Elee) GET(path string, handler HandlerFunc){
	elee.addRouter(http.MethodGet,path,handler)
}

func (elee *Elee) HEAD (path string, handler HandlerFunc){
	elee.addRouter(http.MethodHead,path,handler)
}
func (elee *Elee) POST (path string, handler HandlerFunc){
	elee.addRouter(http.MethodPost,path,handler)
}
func (elee *Elee) PUT(path string, handler HandlerFunc){
	elee.addRouter(http.MethodPut,path,handler) }
func (elee *Elee) PATCH(path string, handler HandlerFunc){
	elee.addRouter(http.MethodPatch,path,handler)
}
func (elee *Elee) DELETE(path string, handler HandlerFunc){
	elee.addRouter(http.MethodDelete,path,handler)
}
func (elee *Elee) CONNECT(path string, handler HandlerFunc){
	elee.addRouter(http.MethodConnect,path,handler)
}
func (elee *Elee) OPTIONS(path string, handler HandlerFunc){
	elee.addRouter(http.MethodOptions,path,handler)
}
func (elee *Elee) TRACE(path string, handler HandlerFunc){
	elee.addRouter(http.MethodTrace,path,handler)
}













func (elee *Elee) ServeHTTP(w http.ResponseWriter, req *http.Request){
	cnt := newContext(w,req)
	meth,ok:= elee.routers[req.URL.Path]
	if !ok{
		cnt.String(http.StatusNotFound,"error:%d\n",http.StatusNotFound)
		return
	}
	handle := meth.findHandler(cnt.Method())
	 if handle == nil{
		 cnt.String(http.StatusNotFound,"error:%d\n",http.StatusNotFound)
		 return
	 }
	 handle(cnt)

}
