package transactions

import "sync"

type ObjectCache struct {
	mu      sync.RWMutex
	objects map[string]map[string]any
	custom  map[string]any
}

func NewObjectCache() *ObjectCache {
	return &ObjectCache{objects: map[string]map[string]any{}, custom: map[string]any{}}
}

func (c *ObjectCache) SetObject(id string, data map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.objects[id] = data
}

func (c *ObjectCache) GetObject(id string) (map[string]any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.objects[id]
	return v, ok
}

func (c *ObjectCache) DeleteObject(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.objects, id)
}

func (c *ObjectCache) ClearOwnedObjects() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.objects = map[string]map[string]any{}
}

func (c *ObjectCache) SetCustom(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.custom[key] = value
}

func (c *ObjectCache) GetCustom(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.custom[key]
	return v, ok
}

func (c *ObjectCache) DeleteCustom(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.custom, key)
}

func (c *ObjectCache) ClearCustom() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.custom = map[string]any{}
}
