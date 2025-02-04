package main

import (
	"container/list"
	"fmt"
)

type Item struct {
	Key   string
	Value interface{}
}

type LRU struct {
	capacity int
	items    map[string]*list.Element
	queue    *list.List
}

func NewLru(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *LRU) Set(key string, value interface{}) bool {
	if element, exists := c.items[key]; exists == true {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.Key] = element

	return true
}

func (c *LRU) purge() {
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*Item)
		delete(c.items, item.Key)
	}
}

func (c *LRU) Get(key string) interface{} {
	element, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.queue.MoveToFront(element)
	return element.Value.(*Item).Value
}

func main() {
	cache := NewLru(5)
	cache.Set("foo1", "bar1")
	cache.Set("foo2", "bar2")
	cache.Set("foo3", "bar3")
	cache.Set("foo4", "bar4")
	cache.Set("foo5", "bar5")
	cache.Set("foo6", "bar6")

	value := cache.Get("foo3")
	vstr, ok := value.(string)
	if ok {
		fmt.Println(vstr)
	}

	front := cache.queue.Front()

	fmt.Printf("front element:%v\n", front.Value)
	n := front.Next()
	for n != nil {
		fmt.Printf("next value:%v ", n.Value)
		n = n.Next()
	}

	fmt.Println(value)

}
