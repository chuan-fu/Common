package db

import (
	"github.com/chuan-fu/Common/db/mysql"
	dbRedis "github.com/chuan-fu/Common/db/redis"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetGorm() *gorm.DB {
	return mysql.GetGorm()
}

func GetRedisCli() redis.Cmdable {
	return dbRedis.GetRedisCli()
}
