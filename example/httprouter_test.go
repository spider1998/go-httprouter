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

func Print11(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "111111111111111111111111!\n")
}

func Print22(w http.ResponseWriter, r *http.Request, ps g.Params) {
	_, _ = fmt.Fprint(w, "22222222222222222222222222!\n")
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
	test1 := router.Group("/test1", 1)
	test1.GET("/index", Index, Print)
	test2 := router.Group("/test2", 2)
	test2.GET("/hello", Hello)

	tes11 := test1.Group("/aaa", 11)
	tes11.GET("/index", Print11)
	tes12 := test1.Group("/bbb", 12)
	tes12.GET("/index", Print22)

	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
