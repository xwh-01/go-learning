package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Db      *gorm.DB
	RedisDB *redis.Client
)
