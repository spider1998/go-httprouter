package go_httprouter

import (
	"golang.org/x/net/context"
	"net/http"
)

type Router struct {
	Trees                  map[string]*Tree
	PanicHandler           func(http.ResponseWriter, *http.Request)
	NotFound               http.Handler
	PathSlash              bool
	UseEscapedPath         bool
	HandleOptions          bool
	HandleMethodNotAllowed bool
	MethodNotAllowed       http.Handler
}

func New() *Router {
	return &Router{
		PathSlash:              true,
		HandleMethodNotAllowed: true,
		HandleOptions:          true,
	}
}

/*METHOD*/
func (r *Router) GET(path string, handles ...Handle) {
	r.Handle("GET", path, handles)
}

func (r *Router) HEAD(path string, handles ...Handle) {
	r.Handle("HEAD", path, handles)
}

func (r *Router) OPTIONS(path string, handles ...Handle) {
	r.Handle("OPTIONS", path, handles)
}

func (r *Router) POST(path string, handles ...Handle) {
	r.Handle("POST", path, handles)
}

func (r *Router) PUT(path string, handles ...Handle) {
	r.Handle("PUT", path, handles)
}

func (r *Router) PATCH(path string, handles ...Handle) {
	r.Handle("PATCH", path, handles)
}

func (r *Router) DELETE(path string, handles ...Handle) {
	r.Handle("DELETE", path, handles)
}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (r *Router) Handle(method, path string, handles []Handle) {
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
			for _, handle := range handles.([]Handle) {
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
		[]Handle{
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
