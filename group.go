package go_httprouter

import (
	"strconv"
)

type RouterGroup struct {
	Groups   map[int]string
	Handlers map[int][]Handler
}

func NewGroup() *RouterGroup {
	return &RouterGroup{
		Groups:   map[int]string{},
		Handlers: map[int][]Handler{},
	}
}

func (r *RouterGroup) Save(index int, handlers ...Handler) {
	r.Handlers[index] = handlers
}

func (r *RouterGroup) Generate(index int) (handlers []Handler) {
	if len(strconv.Itoa(index)) == 1 {
		return r.Handlers[index]
	} else {
		for i := 0; i <= len(strconv.Itoa(index)); i++ {
			k := strconv.Itoa(index)[:i]
			key, _ := strconv.Atoi(k)
			handlers = append(handlers, r.Handlers[key]...)
		}
	}
	return
}

func (r *RouterGroup) Insert(index int, prefix string) {
	r.Groups[index] = prefix
}

func (r *RouterGroup) Get(index int) (path string) {
	if len(strconv.Itoa(index)) == 1 {
		return r.Groups[index]
	} else {
		for i := 0; i <= len(strconv.Itoa(index)); i++ {
			k := strconv.Itoa(index)[:i]
			key, _ := strconv.Atoi(k)
			path += r.Groups[key]
		}
	}
	return
}
