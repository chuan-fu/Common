package redis

import (
	"context"

	"github.com/chuan-fu/Common/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisConf struct {
	Addr     string `required:"true"`
	Password string
	DB       int `default:"0"`
}

var redisCli *redis.Client

func GetRedisCli() *redis.Client {
	return redisCli
}

func ReloadRedis(conf RedisConf) error {
	return ConnectRedis(conf)
}

func ConnectRedis(conf RedisConf) error {
	rc := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	pong, err := rc.Ping(context.TODO()).Result()
	if err != nil {
		err = errors.Wrap(err, "redis连接错误")
		log.Error(err)
		return err
	}
	log.Info("redis pong:", pong)

	redisCli = rc
	return nil
}