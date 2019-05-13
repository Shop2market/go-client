package cache

import (
	"errors"
	"time"
)

const CACHE_TTL = 1 * time.Hour

type Cache struct {
	Data map[string][][]string
	Date *time.Time
}

func New(cached map[string][][]string) *Cache {
	now := time.Now().UTC()
	return &Cache{cached, &now}
}

func (c *Cache) Get() (mappings map[string][][]string, err error) {
	mappings = c.Data
	if c.IsValid() {
		return
	}
	return nil, errors.New("cache is not valid!")
}

func (c *Cache) Update(mappings map[string][][]string) {
	now := time.Now().UTC()
	c.Data = mappings
	c.Date = &now
}

func (c *Cache) IsValid() bool {
	if c.IsEmpty() {
		return false
	}
	if c.IsOutdated() {
		return false
	}

	return true
}

func (c *Cache) IsOutdated() bool {
	now := time.Now().UTC()
	expiration := c.Date.Add(CACHE_TTL)
	dObj := &expiration
	return dObj.Before(now)
}

func (c *Cache) IsEmpty() bool {
	return c.Data == nil || c.Date == nil
}
