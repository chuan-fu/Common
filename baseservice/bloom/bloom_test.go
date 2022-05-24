package bloom

import (
	"fmt"
	"testing"

	"github.com/chuan-fu/Common/db"
	"github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/zlog"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: "127.0.0.1:6379",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	b := NewBloomFilter(db.GetRedisCli(), "bloom:test:1", 1000)
	fmt.Println(b.AddStr("Add1"))
	fmt.Println(b.AddStr("Add501"))
	fmt.Println(b.AddStr("Add503"))
	//for i := 0; i < 500; i++ {
	//	fmt.Println(i, b.AddStr(fmt.Sprintf("Add%d", i)))
	//}
	for i := 480; i < 520; i++ {
		fmt.Print(i)
		fmt.Println(b.ExistsStr(fmt.Sprintf("Add%d", i)))
	}
}