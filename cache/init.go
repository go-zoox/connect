package cache

import "time"

var c = New()

func Get(key string, value interface{}) error {
	return c.core.Get(key, value)
}

func Set(key string, value interface{}, ttl time.Duration) error {
	return c.core.Set(key, value, ttl)
}

func Del(key string) error {
	return c.core.Delete(key)
}
