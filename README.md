# go-httprouter
High-performance HTTP request router based on Go language


Usage

package mian

import (
	`fmt`
	`log`
	`net/http`
	route `github.com/spider1998/go-httprouter`
)

func Index(w http.ResponseWriter, r *http.Request, ps route.Params) {
	_, _ = fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps route.Params) {
	_, _ = fmt.Fprint(w, "Hello World!\n")
}

func main() {
	router := route.New()
	router.GET("/", Index)
	router.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
