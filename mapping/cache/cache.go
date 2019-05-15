package cache

import (
	"errors"
	"time"
)

const CACHE_TTL = 1 * time.Hour

type Cache struct {
	data map[string][][]string
	date *time.Time
}

func New(cached map[string][][]string) *Cache {
	now := time.Now().UTC()
	return &Cache{cached, &now}
}

func NewWithTime(cached map[string][][]string, t time.Time) *Cache {
	return &Cache{cached, &t}
}

func (c *Cache) Get() (mappings map[string][][]string, err error) {
	mappings = c.data
	if c.IsValid() {
		return
	}
	return nil, errors.New("cache is not valid!")
}

func (c *Cache) Update(mappings map[string][][]string) {
	now := time.Now().UTC()
	c.data = mappings
	c.date = &now
}

func (c *Cache) IsValid() bool {
	return !c.IsEmpty() && !c.IsOutdated()
}

func (c *Cache) IsOutdated() bool {
	now := time.Now().UTC()
	if c.data == nil {
		return true
	}
	if c.date != nil {
		expiration := c.date.Add(CACHE_TTL)
		dObj := &expiration
		return dObj.Before(now)
	}
	return true
}

func (c *Cache) IsEmpty() bool {
	return c.data == nil || c.date == nil
}
