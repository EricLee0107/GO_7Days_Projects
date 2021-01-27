package elee

import (
	"net/http"
)



type Elee struct{
	// 静态路由
	routers map[string] *methodHandler
}


func New() *Elee{
	elee := &Elee{
		routers: make(map[string] *methodHandler),
	}
	return elee
}

func (elee *Elee) AddRouter(method, path string, handler HandlerFunc){
	meth,ok := elee.routers[path]
	if !ok{
		meth = &methodHandler{}
		elee.routers[path] = meth
	}
	meth.addHandler(method,handler)
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
