package main

import (
	"fmt"
	g "go-httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "Hello World!\n")
}

func main() {
	router := g.New()
	router.GET("/", Index)
	router.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
