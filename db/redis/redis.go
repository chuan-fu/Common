package redis

import (
	"context"

	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

/*
type Config struct {
	Redis redis.RedisConf `json:"redis" yaml:"redis"`
}

redis:
	addr: 127.0.0.1:6379
	password:
	db: 0

*/

type RedisConf struct {
	Addr     string `required:"true" json:"addr" yaml:"addr"`
	Password string `default:"" json:"password" yaml:"password"`
	DB       int    `default:"0" json:"db" yaml:"db"`
}

var redisCli redis.Cmdable

func GetRedisCli() redis.Cmdable {
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

const (
	Pong = "PONG"
)

func Ping(store redis.Cmdable) bool {
	v, err := store.Ping(context.TODO()).Result()
	if err != nil {
		log.Error(err)
		return false
	}
	return v == Pong
}
