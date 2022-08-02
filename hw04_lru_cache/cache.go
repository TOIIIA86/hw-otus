package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Len() int
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, exists := c.items[key]
	if exists {
		item.Value = newCacheItem(key, value)
		c.queue.MoveToFront(item)
	}
	if !exists {
		item := c.queue.PushFront(newCacheItem(key, value))
		c.items[key] = item
	}
	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		delete(c.items, back.Value.(*cacheItem).key)
		c.queue.Remove(back)
	}
	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, exists := c.items[key]
	if exists {
		ci := item.Value.(*cacheItem)
		c.queue.MoveToFront(item)
		return ci.value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Len() int {
	return c.queue.Len()
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	ci := new(cacheItem)
	ci.key = key
	ci.value = value
	return ci
}
