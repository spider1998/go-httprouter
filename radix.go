//Package gohttprouter 提供路由处理存储
package gohttprouter

import (
	"sort"
	"strings"
)

//Tree 存储实体结构
type Tree struct {
	Root *Node
	Size int
}

//Node 节点
type Node struct {
	Leaf   *LeafNode
	Prefix string
	Sons   Sons
}

//LeafNode 叶子节点
type LeafNode struct {
	Key   string
	Value interface{}
}

//NewTree 创建树
func NewTree() *Tree {
	return &Tree{
		Root: &Node{},
	}
}

//Insert 存储路由处理信息
//	t.Insert(path, handles)
//基数树
func (t *Tree) Insert(key string, value interface{}) (interface{}, bool) {
	search := key
	n := t.Root
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				oldValue := n.Leaf.Value
				n.Leaf.Value = value
				return oldValue, true
			}
			n.Leaf = &LeafNode{
				Key:   key,
				Value: value,
			}
			t.Size++
			return nil, false
		}

		parent := n
		n = n.getSon(search[0])
		if n == nil {
			s := Son{
				Label: search[0],
				Node: &Node{
					Leaf: &LeafNode{
						Key:   key,
						Value: value,
					},
					Prefix: search,
				},
			}
			parent.addSon(s)
			t.Size++
			return nil, false

		}

		index := 0
		for ; index < len(search) && index < len(n.Prefix); index++ {
			if search[index] != n.Prefix[index] {
				break
			}
		}
		if index == len(n.Prefix) {
			search = search[index:]
			continue
		}

		t.Size++
		child := &Node{
			Prefix: search[:index],
		}
		parent.updateSon(search[0], child)

		// Restore the existing node
		child.addSon(Son{
			Label: n.Prefix[index],
			Node:  n,
		})
		n.Prefix = n.Prefix[index:]

		leaf := &LeafNode{
			Key:   key,
			Value: value,
		}

		search = search[index:]
		if len(search) == 0 {
			child.Leaf = leaf
			return nil, false
		}

		child.addSon(Son{
			Label: search[0],
			Node: &Node{
				Leaf:   leaf,
				Prefix: search,
			},
		})
		return nil, false
	}
}

//Get 通过路由获取处理函数
//	if handles, _ := root.Get(path); handles != nil {
//			for _, handle := range handles.([]Handle) {
//				handle(w, req, nil)
//			}
//			return
//		}
//返回接口类型
func (t *Tree) Get(s string) (interface{}, bool) {
	n := t.Root
	search := s
	for {
		// Check for key exhaution
		if len(search) == 0 {
			if n.isLeaf() {
				return n.Leaf.Value, true
			}
			break
		}

		// Look for an edge
		n = n.getSon(search[0])
		if n == nil {
			break
		}

		// Consume the search prefix
		if strings.HasPrefix(search, n.Prefix) {
			search = search[len(n.Prefix):]
		} else {
			break
		}
	}
	return nil, false
}

func (n *Node) isLeaf() bool {
	return n.Leaf != nil
}

func (n *Node) getSon(label byte) *Node {
	num := len(n.Sons)
	idx := sort.Search(num, func(i int) bool {
		return n.Sons[i].Label >= label
	})
	if idx < num && n.Sons[idx].Label == label {
		return n.Sons[idx].Node //返回新的节点
	}
	return nil
}

func (n *Node) addSon(s Son) {
	n.Sons = append(n.Sons, s)
	n.Sons.Sort()
}

func (n *Node) updateSon(label byte, node *Node) {
	num := len(n.Sons)
	idx := sort.Search(num, func(i int) bool {
		return n.Sons[i].Label >= label
	})
	if idx < num && n.Sons[idx].Label == label {
		n.Sons[idx].Node = node
		return
	}
	panic("replacing missing son")
}

//Sons 子节点集合
type Sons []Son

//Son 子节点
type Son struct {
	Label byte
	Node  *Node
}

//Len 获取长度
func (s Sons) Len() int {
	return len(s)
}

//Less 比较大小
func (s Sons) Less(i, j int) bool {
	return s[i].Label < s[j].Label
}

//Swap 交换值
func (s Sons) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Sort 排序
func (s Sons) Sort() {
	sort.Sort(s)
}

//AddRoute 添加路由
func (t *Tree) AddRoute(path string, handlers []Handler) {
	t.Insert(path, handlers)
}
