package midleware

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	cleanupTime int64         //время очистки
	expiration  time.Duration //продолжительность жизни
	items       map[string]float64
	mu          sync.RWMutex
}

func New(expiration time.Duration) *Cache {
	items := make(map[string]float64)

	if expiration <= 0 {
		return nil
	}

	expire := time.Now().Add(expiration)
	cache := Cache{
		items:       items,
		expiration:  expiration,
		cleanupTime: expire.UnixNano(),
	}

	go cache.StartGC()

	return &cache
}

func (c *Cache) Count() int {
	return len(c.items)
}

func (c *Cache) Expired() bool {
	return time.Now().UnixNano() > c.cleanupTime
}

func (c *Cache) BulkSet(rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = rates
	expire := time.Now().Add(c.expiration)
	c.cleanupTime = expire.UnixNano()
}

func (c *Cache) Set(key string, value float64) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache) GetAll() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := map[string]float64{}
	for k, v := range c.items {
		result[k] = v
	}
	return result
}

func (c *Cache) Get(key string, value float64) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().UnixNano() > c.cleanupTime {
		c.mu.RLock()
		return 0, false
	}

	item, ok := c.items[key]
	return item, ok
}
func (c *Cache) Return(key string) map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *Cache) StartGC() {
	for {
		<-time.After(c.expiration)

		if c.items == nil {
			return
		}

		if time.Now().UnixNano() > c.cleanupTime {
			c.mu.Lock()
			// c.items = map[string]float64{}
			clear(c.items)
			c.mu.Unlock()
		}
	}
}
