package cache

import (
	"sync"
)

type Cache interface {
	Get(key string) interface{}
	Add(key string, value interface{})
	Update(key string, newValue interface{})
	Delete(key string)
	Exists(key string) bool
	Values() []interface{}
	Keys() []string
	ToMap() map[string]interface{}
}

type cache struct {
	cacheLock sync.RWMutex
	cache     map[string]interface{}
}

func (c *cache) toMapInternal() map[string]interface{} {
	kvMap := make(map[string]interface{})
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	for key, value := range c.cache {
		kvMap[key] = value
	}
	return kvMap
}

func (c *cache) ToMap() map[string]interface{} {
	return c.toMapInternal()
}

func (c *cache) getKeysInternal() []string {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	var keys []string
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys
}

func (c *cache) Keys() []string {
	return c.getKeysInternal()
}

func (c *cache) Exists(key string) bool {
	return c.Get(key) != nil
}

func (c *cache) Add(key string, value interface{}) {
	c.addInternal(key, value)
}

func (c *cache) addInternal(key string, newValue interface{}) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	if _, exists := c.cache[key]; !exists {
		c.cache[key] = newValue
	}
}

func (c *cache) updateInternal(key string, newValue interface{}) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	if _, exists := c.cache[key]; exists {
		c.cache[key] = newValue
	}
}

func (c *cache) getInternal(key string) interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	if value, exists := c.cache[key]; exists {
		return value
	}
	return nil
}

func (c *cache) deleteInternal(key string) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	delete(c.cache, key)
}

func (c *cache) getValuesInternal() []interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	var values []interface{}
	for _, value := range c.cache {
		values = append(values, value)
	}
	return values
}

func (c *cache) Get(key string) interface{} {
	return c.getInternal(key)
}

func (c *cache) Update(key string, newValue interface{}) {
	c.updateInternal(key, newValue)
}

func (c *cache) Delete(key string) {
	c.deleteInternal(key)
}

func (c *cache) Values() []interface{} {
	return c.getValuesInternal()
}

func Default() Cache {
	return &cache{
		cacheLock: sync.RWMutex{},
		cache:     make(map[string]interface{}),
	}
}
