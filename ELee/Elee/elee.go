package Elee

import (
	"fmt"
	"net/http"
)



type Elee struct{
	routers map[string] http.HandlerFunc
}

func New() *Elee{
	elee := &Elee{
		routers: make(map[string] http.HandlerFunc),
	}
	return elee
}

func (elee *Elee) AddRouter(path string, handler http.HandlerFunc){
	elee.routers[path] = handler
}



func (elee *Elee) ServeHTTP(w http.ResponseWriter, req *http.Request){
	handle,ok:= elee.routers[req.URL.Path]
	if !ok{
		fmt.Fprintf(w,"error:%d\n",http.StatusNotFound)
		return
	}
	handle(w,req)
}
