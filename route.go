//router.go文件

//提供简单http路由
package go_httprouter

import (
	"golang.org/x/net/context"
	"net/http"
)

//Router 路由结构体
type Router struct {
	RouteGroup
	Trees                  map[string]*Tree
	PanicHandler           func(http.ResponseWriter, *http.Request)
	NotFound               http.Handler
	PathSlash              bool
	UseEscapedPath         bool
	HandleOptions          bool
	HandleMethodNotAllowed bool
	MethodNotAllowed       http.Handler
}

type RouteGroup struct {
	Prefix   string
	Handlers []Handler
	Router   *Router
}

//New 创建新的路由
//	router := g.New()
//返回路由实例
func New() *Router {
	return &Router{
		PathSlash:              true,
		HandleMethodNotAllowed: true,
		HandleOptions:          true,
	}
}

func (r *Router) Group(prefix string, handlers []Handler) *Router {
	r.RouteGroup = RouteGroup{
		Prefix:   prefix,
		Handlers: handlers,
	}
	return r
}

//GET get方法
//	router.GET("/", Index)
//以下方法类似
func (r *Router) GET(path string, handlers ...Handler) {
	r.Handle("GET", r.Prefix+path, combineHandlers(r, handlers))
}

//HEAD head方法
func (r *Router) HEAD(path string, handlers ...Handler) {
	r.Handle("HEAD", r.Prefix+path, combineHandlers(r, handlers))
}

//OPTIONS options方法
func (r *Router) OPTIONS(path string, handlers ...Handler) {
	r.Handle("OPTIONS", r.Prefix+path, combineHandlers(r, handlers))
}

//POST post方法
func (r *Router) POST(path string, handlers ...Handler) {
	r.Handle("POST", r.Prefix+path, combineHandlers(r, handlers))
}

//PUT put方法
func (r *Router) PUT(path string, handlers ...Handler) {
	r.Handle("PUT", r.Prefix+path, combineHandlers(r, handlers))
}

//PATCH patch方法
func (r *Router) PATCH(path string, handlers ...Handler) {
	r.Handle("PATCH", r.Prefix+path, combineHandlers(r, handlers))
}

//DELETE delete方法
func (r *Router) DELETE(path string, handlers ...Handler) {
	r.Handle("DELETE", r.Prefix+path, combineHandlers(r, handlers))
}

//Handle 路由处理函数
type Handler func(http.ResponseWriter, *http.Request, Params)

//Param 路由处理参数（暂留）
type Param struct {
	Key   string
	Value string
}

type Params []Param

//Handle 处理路由函数
//	r.Handle("GET", path, handles)
//存储路由
func (r *Router) Handle(method, path string, handles []Handler) {
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

	root.AddRoute(path, handles)
}

//ServeHTTP 实现ServeHTTP方法
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	//Handle_Panic
	if r.PanicHandler != nil {
		defer r.handlePanic(w, req)
	}

	//USE_SLASH
	if r.PathSlash {
		if len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}

	//USE_ESCAPES
	if r.UseEscapedPath {
		path = req.URL.EscapedPath()
	}

	//Handle
	if root := r.Trees[req.Method]; root != nil {
		if handles, _ := root.Get(path); handles != nil {
			for _, handle := range handles.([]Handler) {
				handle(w, req, nil)
			}
			return
		}
	}

	//OPTIONS and METHOD_NOT_ALLOWED
	if r.HandleOptions && req.Method == "OPTIONS" {
		if allow := r.allowedMethod(path, req.Method); allow != "" {
			w.Header().Set("Allow", allow)
			return
		}
	} else {
		if r.HandleMethodNotAllowed {
			if allow := r.allowedMethod(path, req.Method); allow != "" {
				w.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

//HandlerFunc 内置处理函数
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
		[]Handler{
			func(w http.ResponseWriter, req *http.Request, p Params) {
				ctx := req.Context()
				ctx = context.WithValue(ctx, ParamsKey, p)
				req = req.WithContext(ctx)
				handler.ServeHTTP(w, req)
			},
		},
	)
}

func (r *Router) allowedMethod(path, method string) (methods string) {
	if path == "*" {
		for key := range r.Trees {
			if key == "OPTIONS" {
				continue
			}
			if methods == "" {
				methods = key
			} else {
				methods += "," + key
			}
		}
	} else {
		for key := range r.Trees {
			if key == method || key == "OPTIONS" {
				continue
			}
			_, exist := r.Trees[key].Get(path)
			if exist {
				if methods == "" {
					methods = key
				} else {
					methods += "," + key
				}
			}
		}
	}
	if len(methods) > 0 {
		methods += "," + "OPTIONS"
	}
	return
}

func (r *Router) handlePanic(w http.ResponseWriter, req *http.Request) {
	if re := recover(); re != nil {
		r.PanicHandler(w, req)
	}
}

func combineHandlers(r *Router, handles []Handler) (handlers []Handler) {
	return append(append(handlers, r.Handlers...), handles...)
}
