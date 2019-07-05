package mian

import (
	"fmt"
	"log"
	"net/http"
	"test_route/radix"
)

func Index(w http.ResponseWriter, r *http.Request, ps radix.Params) {
	_, _ = fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps radix.Params) {
	_, _ = fmt.Fprint(w, "Hello World!\n")
}

func main() {
	router := radix.New()
	router.GET("/", Index)
	router.GET("/hello", Hello)
	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
