package elee

import (
	"fmt"
	"net/http"
)

// 同一路径不同请求的处理器集合
type methodHandler struct {
		connect  HandlerFunc
		delete   HandlerFunc
		get      HandlerFunc
		head     HandlerFunc
		options  HandlerFunc
		patch    HandlerFunc
		post     HandlerFunc
		put      HandlerFunc
		trace    HandlerFunc
	}


func (mh *methodHandler) addHandler (method string, h HandlerFunc){
	switch method {
	case http.MethodConnect:
		mh.connect = h
	case http.MethodDelete:
		mh.delete = h
	case http.MethodGet:
		mh.get = h
	case http.MethodHead:
		mh.head = h
	case http.MethodOptions:
		mh.options = h
	case http.MethodPatch:
		mh.patch = h
	case http.MethodPost:
		mh.post = h
	case http.MethodPut:
		mh.put = h
	case http.MethodTrace:
		mh.trace = h
	default:
		fmt.Println("error method:",method)
	}
}

func (mh *methodHandler) findHandler (method string) HandlerFunc{
	switch method {
	case http.MethodConnect:
		return mh.connect
	case http.MethodDelete:
		return mh.delete
	case http.MethodGet:
		return mh.get
	case http.MethodHead:
		return mh.head
	case http.MethodOptions:
		return mh.options
	case http.MethodPatch:
		return mh.patch
	case http.MethodPost:
		return mh.post
	case http.MethodPut:
		return mh.put
	case http.MethodTrace:
		return mh.trace
	default:
		return nil
	}

}

