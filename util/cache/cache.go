package cache

import (
	"minik8s/apiObject/types"
	"sync"
)

type Cache interface {
	Get(podUID types.UID) interface{}
	Update(podUID types.UID, newValue interface{})
	Delete(podUID types.UID)
	Values() []interface{}
}

type cache struct {
	cacheLock sync.RWMutex
	cache     map[types.UID]interface{}
}

func (c *cache) updateInternal(podUID types.UID, newValue interface{}) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.cache[podUID] = newValue
}

func (c *cache) getInternal(podUID types.UID) interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	if value, exists := c.cache[podUID]; exists {
		return value
	}
	return nil
}

func (c *cache) deleteInternal(podUID types.UID) {
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

func (c *cache) Get(podUID types.UID) interface{} {
	return c.getInternal(podUID)
}

func (c *cache) Update(podUID types.UID, newValue interface{}) {
	c.updateInternal(podUID, newValue)
}

func (c *cache) Delete(podUID types.UID) {
	c.deleteInternal(podUID)
}

func (c *cache) Values() []interface{} {
	return c.getValuesInternal()
}

func Default() Cache {
	return &cache{
		cacheLock: sync.RWMutex{},
		cache:     make(map[types.UID]interface{}),
	}
}
