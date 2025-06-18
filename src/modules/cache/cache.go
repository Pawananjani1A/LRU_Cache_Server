/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 03 Apr 2024
*/

package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key        string
	value      interface{}
	expiration time.Time
}

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	lruList  *list.List
	mutex    sync.Mutex
}

var instance *LRUCache
var once sync.Once

func GetLRUCache(capacity int) *LRUCache {
	once.Do(func() {
		instance = &LRUCache{
			capacity: capacity,
			cache:    make(map[string]*list.Element),
			lruList:  list.New(),
		}
	})
	return instance
}

//func NewLRUCache(capacity int) *LRUCache {
//	return &LRUCache{
//		capacity: capacity,
//		cache:    make(map[string]*list.Element),
//		lruList:  list.New(),
//	}
//}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		if elem.Value.(*entry).expiration.Before(time.Now()) {
			c.removeElement(elem)
			return nil, false
		}
		c.lruList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.lruList.MoveToFront(elem)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiration = time.Now().Add(expiration)
	} else {
		if len(c.cache) >= c.capacity {
			c.evictOldest()
		}
		elem := c.lruList.PushFront(&entry{key, value, time.Now().Add(expiration)})
		c.cache[key] = elem
	}
}

func (c *LRUCache) evictOldest() {
	elem := c.lruList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.lruList.Remove(e)
	delete(c.cache, e.Value.(*entry).key)
}
