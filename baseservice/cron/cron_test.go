package cron

import (
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/mutex"
	"github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/zlog"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestRunCronTask(t *testing.T) {
	mutex.DefaultExpiration = time.Second * 2
	c1 := NewCronTask(
		"* * * * *",
		func() {
			log.Info("run * * * * *")
		},
		WithCondition(func() bool { return true }),
		WithMutex("key1"),
	)

	c2 := NewCronTask(
		"* * * * *",
		func() {
			log.Info("always run * * * * *")
		},
	)

	c3 := NewCronTask(
		"* * * * *",
		func() {
			log.Info("always not run * * * * *")
		},
		WithCondition(func() bool { return false }),
	)

	RunCronTask(
		c1,
		c2,
		c3,
	)

	time.Sleep(time.Minute)
}
