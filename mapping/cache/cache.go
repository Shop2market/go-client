package cache

import (
	"errors"
	"time"
)

const CACHE_TTL = 1 * time.Hour

type Cache struct {
	mappings map[string][][]string
	date     time.Time
}

func New() *Cache {
	return &Cache{map[string][][]string{}, time.Now().UTC()}
}

func NewWithTime(cached map[string][][]string, t time.Time) *Cache {
	return &Cache{cached, t}
}

func (c *Cache) Get() (mappings map[string][][]string, err error) {
	mappings = c.mappings
	if c.IsValid() {
		return
	}
	return nil, errors.New("cache is not valid!")
}

func (c *Cache) Update(mappings map[string][][]string) {
	c.mappings = mappings
	c.date = time.Now().UTC()
}

func (c *Cache) IsValid() bool {
	if len(c.mappings) == 0 {
		return false
	}

	expiration := c.date.Add(CACHE_TTL)
	dObj := &expiration
	return dObj.After(time.Now().UTC())
}
