package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var client *redis.Client

func InitLike() {
	client = redis.NewClient(&redis.Options{
		Addr:        "newk8s.ferdinandaedth.top:6379",
		Password:    "123456",         // Redis数据库没有密码
		DB:          0,                // 默认数据库为0
		PoolSize:    10,               // 设置连接池大小
		PoolTimeout: 10 * time.Second, // 连接池等待超时时间
	})
}
