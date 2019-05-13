package cache

import "time"

const CACHE_TTL = 1 * time.Hour

type Cache struct {
	Data map[string][][]string
	Date *time.Time
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
	now := time.Now().UTC()
	expiration := c.Date.Add(CACHE_TTL)
	return expiration.After(now)
}

func (c *Cache) IsEmpty() bool {
	return c.Data == nil || c.Date == nil
}
