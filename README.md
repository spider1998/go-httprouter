# go-httprouter
[![Build Status](https://travis-ci.org/spider1998/go-httprouter.svg?branch=master)](https://travis-ci.org/spider1998/go-httprouter)
[![GoDoc](http://godoc.org/github.com/spider1998/go-httprouter?status.svg)](http://godoc.org/github.com/spider1998/go-httprouter)
![](https://img.shields.io/badge/language-Go-orange.svg)


轻量级高性能HTTP路由框架（持续更新）

High-performance HTTP request router based on Go language

![](https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1562330784049&di=bba8782630e41c2b0399600e48a1f9e2&imgtype=0&src=http%3A%2F%2Fimg.mp.itc.cn%2Fupload%2F20161129%2F130444cd837c49c7bef4239afe39dc2f.jpg)  


Usage
```Go
package main

import (
	"fmt"
	g "github.com/spider1998/go-httprouter"
	`log`
	"net/http"
	`testing`
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

func TestRouterGroup(t *testing.T) {
	router := g.New()
	test1 := router.Group("/test1",1)
	test1.GET("/index", Index,Print)
	test2 := router.Group("/test2",2)
	test2.GET("/hello", Hello)

	tes11:= test1.Group("/aaa",11)
	tes11.GET("/index",Print11)
	tes12:= test1.Group("/bbb",12)
	tes12.GET("/index",Print22)


	log.Println("start run...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
```
