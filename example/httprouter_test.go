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

func Print(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "++++++++++++++++Hello World!\n")
}

/*func TestHttptouter(t *testing.T) {
	router := g.New()
	router.Group("/test1",nil)
	router.GET("/", Index)
	router.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}*/

func TestRouterGroup(t *testing.T) {
	router := g.New()
	test1 := router.Group("/test1", nil)
	test1.GET("/index", Index, Print)
	test2 := router.Group("/test2", nil)
	test2.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
