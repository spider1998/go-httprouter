package go_httprouter

import (
	"strconv"
)

type RouterGroup struct {
	Groups map[int]string
}

func NewGroup() *RouterGroup {
	return &RouterGroup{
		Groups: map[int]string{},
	}
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
