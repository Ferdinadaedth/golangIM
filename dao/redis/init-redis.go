package redis

import "github.com/go-redis/redis"

var client *redis.Client

func InitLike() {
	client = redis.NewClient(&redis.Options{
		Addr:     "newk8s.ferdinandaedth.top:6379",
		Password: "123456", // Redis数据库没有密码
		DB:       0,        // 默认数据库为0
	})
}
