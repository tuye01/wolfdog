package consts

import (
	"github.com/go-redis/redis/v8"
	gorm "gorm.io/gorm"
)

var GormDB *gorm.DB

var RedisSuf = "gosso.cache.redis."
var RedisDB *redis.Client
