package elee

import (
	"net/http"
)



type Elee struct{
	// 静态路由
	//routers map[string] *methodHandler
	tree *node
	maxParam *int
}


func New() *Elee{
	elee := &Elee{
		maxParam:        new(int),

		//routers: make(map[string] *methodHandler),
		tree: &node{
			// 节点对应的不同http metthod 的handler
			methodHandler: new(methodHandler),
		},
	}
	elee.tree.elee= elee
	return elee
}

func (elee *Elee) addRouter(method, path string, handler HandlerFunc){
	elee.tree.Add(method,path,handler)
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
	elee.tree.Find(req.Method,req.URL.Path,cnt)
	h := cnt.Handler()
	if h == nil{
		cnt.String(http.StatusNotFound,"error:%d\n",http.StatusNotFound)
		return
	}
	h(cnt)
}
