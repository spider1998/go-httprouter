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

//加入路由组处理方法
func (r *RouterGroup) HandlerInsert(index int, handlers ...Handler) {
	r.Handlers[index] = handlers
}

//获取路由组所有指定执行的方法
func (r *RouterGroup) HandlerGenerate(index int) (handlers []Handler) {
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

//加入路由组前缀
func (r *RouterGroup) PrefixInsert(index int, prefix string) {
	r.Groups[index] = prefix
}

//生成路由组前缀
func (r *RouterGroup) PrefixGenerate(index int) (path string) {
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
