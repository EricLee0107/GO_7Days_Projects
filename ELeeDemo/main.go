package main

import (
	"elee"
	"fmt"
	"net/http"
)

func main() {
	el := elee.New()
	el.AddRouter(http.MethodGet,"/",indexHandler)
	el.AddRouter(http.MethodPost,"/",indexHandlerPost)
	http.ListenAndServe("localhost:8888",el)
}




func indexHandlerPost(cnt elee.Context) error{
	cnt.String(http.StatusOK,"post this %s page;","index")
	return nil
}
func indexHandler(cnt elee.Context) error{
	cnt.String(http.StatusOK,"this %s page;","index")
	return nil
}

func helloELee(w http.ResponseWriter, req *http.Request){
	fmt.Fprintln(w,"this hello lzw handler")
}

