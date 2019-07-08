# go-httprouter
[![Build Status](https://travis-ci.org/spider1998/go-httprouter.svg?branch=master)](https://travis-ci.org/spider1998/go-httprouter)
[![GoDoc](http://godoc.org/github.com/spider1998/go-httprouter?status.svg)](http://godoc.org/github.com/spider1998/go-httprouter)


轻量级高性能HTTP路由框架

High-performance HTTP request router based on Go language

![](https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1562330784049&di=bba8782630e41c2b0399600e48a1f9e2&imgtype=0&src=http%3A%2F%2Fimg.mp.itc.cn%2Fupload%2F20161129%2F130444cd837c49c7bef4239afe39dc2f.jpg)  


Usage
```Go
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
``
