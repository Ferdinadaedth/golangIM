package main

import (
	"golangIM/api"
	"golangIM/cache"
	"golangIM/dao"
	"golangIM/utils"
)

func main() {
	utils.InitMq()
	dao.InitMysql()
	cache.InitCache()
	api.InitRouter()
}
