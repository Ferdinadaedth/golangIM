package cache

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

var Cache *cache.Cache
var ring *redis.Ring

func InitCache() {
	ring = redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": "newk8s.ferdinandaedth.top:6379",
		},
		Password:   "123456",
		PoolSize:   10, // 设置连接池大小
		MaxRetries: 3,  // 设置最大重试次数
	})
	Cache = cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

}

func SetCache(key string, obj interface{}) error {
	ctx := context.TODO()
	if err := Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}
	return nil
}

func GetCache(key string, wanted interface{}) error {

	ctx := context.TODO()
	if err := Cache.Get(ctx, key, wanted); err != nil {
		return err
	}

	return nil

}

func DeleteCache(key string) error {
	ctx := context.TODO()

	err := Cache.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
