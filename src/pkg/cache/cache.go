/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package cache

import (
	"sync"
	"time"
)

// CacheEntry is a value stored in the cache.
type CacheEntry struct {
	value      interface{}
	expiration int64
	persistent bool // added field to indicate if entry is persistent
}

// SafeCache is a thread-safe cache.
type SafeCache struct {
	syncMap sync.Map
}

// Set stores a value in the cache with a given TTL (time to live) in seconds.
// If ttl is zero, the entry is considered persistent (non-expiring).
func (sc *SafeCache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64
	var persistent bool

	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	} else {
		persistent = true // mark as persistent if ttl is zero
	}

	sc.syncMap.Store(key, CacheEntry{value: value, expiration: expiration, persistent: persistent})
}

// Get retrieves a value from the cache. If the value is not found
// or has expired (and is not persistent), it returns false.
func (sc *SafeCache) Get(key string) (interface{}, bool) {
	entry, found := sc.syncMap.Load(key)
	if !found {
		return nil, false
	}

	cacheEntry := entry.(CacheEntry)
	if !cacheEntry.persistent && time.Now().UnixNano() > cacheEntry.expiration {
		sc.syncMap.Delete(key)
		return nil, false
	}
	return cacheEntry.value, true
}

// Delete removes a value from the cache.
func (sc *SafeCache) Delete(key string) {
	sc.syncMap.Delete(key)
}

// CleanUp periodically removes expired entries from the cache.
func (sc *SafeCache) CleanUp() {
	for {
		time.Sleep(1 * time.Minute)
		sc.syncMap.Range(func(key, entry interface{}) bool {
			cacheEntry := entry.(CacheEntry)
			if !cacheEntry.persistent && time.Now().UnixNano() > cacheEntry.expiration {
				sc.syncMap.Delete(key)
			}
			return true
		})
	}
}
