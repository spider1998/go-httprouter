package go_httprouter

import (
	"golang.org/x/net/context"
	"net/http"
)

type Router struct {
	Trees          map[string]*Tree
	NotFound       http.Handler
	PathSlash      bool
	UseEscapedPath bool
}

func New() *Router {
	return &Router{
		PathSlash: true,
	}
}

/*METHOD*/
func (r *Router) GET(path string, handle Handle) {
	r.Handle("GET", path, handle)
}

func (r *Router) HEAD(path string, handle Handle) {
	r.Handle("HEAD", path, handle)
}

func (r *Router) OPTIONS(path string, handle Handle) {
	r.Handle("OPTIONS", path, handle)
}

func (r *Router) POST(path string, handle Handle) {
	r.Handle("POST", path, handle)
}

func (r *Router) PUT(path string, handle Handle) {
	r.Handle("PUT", path, handle)
}

func (r *Router) PATCH(path string, handle Handle) {
	r.Handle("PATCH", path, handle)
}

func (r *Router) DELETE(path string, handle Handle) {
	r.Handle("DELETE", path, handle)
}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (r *Router) Handle(method, path string, handle Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.Trees == nil {
		r.Trees = make(map[string]*Tree)
	}

	root := r.Trees[method]
	if root == nil {
		root = NewTree()
		r.Trees[method] = root
	}

	root.AddRoute(path, handle)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if r.PathSlash {
		if len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	if r.UseEscapedPath {
		path = req.URL.EscapedPath()
	}
	if root := r.Trees[req.Method]; root != nil {
		if handle, _ := root.Get(path); handle != nil {
			handle.(Handle)(w, req, nil)

			return
		}
	}
	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle. With go 1.7+, the Params will be available in the
// request context under ParamsKey.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, p Params) {
			ctx := req.Context()
			ctx = context.WithValue(ctx, ParamsKey, p)
			req = req.WithContext(ctx)
			handler.ServeHTTP(w, req)
		},
	)
}
