package cache

import (
	"sync"
)

type Cache interface {
	Get(podUID string) interface{}
	Update(podUID string, newValue interface{})
	Delete(podUID string)
	Values() []interface{}
}

type cache struct {
	cacheLock sync.RWMutex
	cache     map[string]interface{}
}

func (c *cache) updateInternal(podUID string, newValue interface{}) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.cache[podUID] = newValue
}

func (c *cache) getInternal(podUID string) interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	if value, exists := c.cache[podUID]; exists {
		return value
	}
	return nil
}

func (c *cache) deleteInternal(podUID string) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	delete(c.cache, podUID)
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

func (c *cache) Get(podUID string) interface{} {
	return c.getInternal(podUID)
}

func (c *cache) Update(podUID string, newValue interface{}) {
	c.updateInternal(podUID, newValue)
}

func (c *cache) Delete(podUID string) {
	c.deleteInternal(podUID)
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
