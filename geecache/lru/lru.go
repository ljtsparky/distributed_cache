package lru

import (
	"container/list"
	// "fmt"
)

type Value interface {
	Len() int // why not int64?
	// 1. int is the natural integer type. Most built-in functions return int
	// 2. familiar pattern: go's builtin len() function returns int
	// 3. practical sizes: most data structures won't need more than int to hold for their length
	// on 64bit systems, int is typically int64. But not for int32. But int64 is always 64 bits.
}

type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List 
	cache    map[string]*list.Element 
	onEvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

func New(maxBytes int64, onEvicted_cb func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:		  list.New(),
		cache:    make(map[string]*list.Element),
		onEvicted: onEvicted_cb,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) 
		// type assertion that ele.Value interface{} is actually a pointer to an entry
		kv := ele.Value.(*entry) 
		// without this unwrap, we can't access kv.value
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Len() int {
	return c.ll.Len() 
}



