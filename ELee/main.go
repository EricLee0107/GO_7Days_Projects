package main

import (
	"ELee/Elee"
	"fmt"
	"net/http"
)

func main() {
	elee := Elee.New()
	elee.AddRouter("/",indexHandler)
	elee.AddRouter("/ELee",helloELee)
	http.ListenAndServe("localhost:8888",elee)
}




func indexHandler(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w,"this index handler\n")
}

func helloELee(w http.ResponseWriter, req *http.Request){
	fmt.Fprintln(w,"this hello lzw handler")
}

