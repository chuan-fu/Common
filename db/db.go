package db

import (
	"context"

	"github.com/chuan-fu/Common/db/mysql"
	dbRedis "github.com/chuan-fu/Common/db/redis"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetGorm(ctx context.Context) *gorm.DB {
	return mysql.GetGorm().WithContext(ctx)
}

func GetGormWithoutCtx() *gorm.DB {
	return mysql.GetGorm()
}

func GetRedisCli() redis.Cmdable {
	return dbRedis.GetRedisCli()
}
