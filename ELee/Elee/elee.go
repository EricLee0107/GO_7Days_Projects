package Elee

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"sync"
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

func (elee *Elee) AddRouter(method, path string, handler http.HandlerFunc){
	meth,ok := elee.routers[path]
	if !ok{
		meth = &methodHandler{}
		elee.routers[path] = meth
	}

}

func ()




func (elee *Elee) ServeHTTP(w http.ResponseWriter, req *http.Request){
	handle,ok:= elee.routers[req.URL.Path]
	if !ok{
		fmt.Fprintf(w,"error:%d\n",http.StatusNotFound)
		return
	}
	handle(w,req)
}
