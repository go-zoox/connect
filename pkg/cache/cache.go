package cache

import (
	"time"

	"github.com/go-zoox/kv"
	"github.com/go-zoox/kv/typing"
)

type Cache struct {
	core kv.KV
}

func New() *Cache {
	core, err := kv.New(&typing.Config{
		Engine: "memory",
	})
	if err != nil {
		panic(err)
	}

	return &Cache{
		core: core,
	}
}

func (c *Cache) Get(key string, value interface{}) error {
	return c.core.Get(key, value)
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	return c.core.Set(key, value, ttl)
}

func (c *Cache) Del(key string) error {
	return c.core.Delete(key)
}
