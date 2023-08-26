package main

import (
	"golangIM/api"
	"golangIM/cache"
	"golangIM/dao"
	"golangIM/dao/redis"
)

func main() {
	dao.Init()
	redis.InitLike()
	cache.InitCache()
	api.InitRouter()
}
