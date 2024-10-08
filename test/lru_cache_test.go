package test

import "github.com/siddontang/go/list2"

type LRUCache struct {
	cache map[int]*Elem
	order *Elem   // 双向链表保存缓存的顺序
	size  int
	head ,tail *Elem
}

type Elem struct {
	key  int
	val  int
	Next,Prev *Elem
}

func Constructor(capacity int) *LRUCache {
	return &LRUCache{
		size:  capacity,
		cache: make(map[int]*Elem, capacity),
	}
}

func (l *LRUCache) Get(key int) int {
	val, ok := l.cache[key]
	if ok {

		l
		return val.val
	}

}

func (l *LRUCache) movetoHead(elem *Elem) int {
	tmp := elem
	prev := elem.Prev
	 = elem.Next

}


func (l *LRUCache) Put(key, value int) int {

}
