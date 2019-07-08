package main

import (
	"fmt"
	g "github.com/spider1998/go-httprouter"
	"log"
	"net/http"
	"testing"
)

func Index(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "Hello World!\n")
}

func TestHttptouter(t *testing.T) {
	router := g.New()
	router.GET("/", Index)
	router.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
