package main

import (
	"golangIM/api"
	"golangIM/cache"
	"golangIM/dao/redis"
)

func main() {
	redis.InitLike()
	cache.InitCache()
	api.InitRouter()
}
