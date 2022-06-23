package tokenlimit

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/chuan-fu/Common/baseservice/cast"

	dbRedis "github.com/chuan-fu/Common/db/redis"

	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	xrate "golang.org/x/time/rate"
)

/*
-- 速率
local rate = tonumber(ARGV[1])
-- 间隔
local per = tonumber(ARGV[2])
-- 桶容量
local cap = tonumber(ARGV[3])
-- 超时时间
local ttl = tonumber(ARGV[4])
-- 当前时间戳
local now = tonumber(ARGV[5])
-- 当前请求token数量
local requested = tonumber(ARGV[6])


-- 当前时间桶容量，为空则代表第一次进入，设置为最大容量
local tokens = tonumber(redis.call("get", KEYS[1]))
if tokens == nil then
	tokens = cap
end

-- 上一次刷新的时间，为空设置为0
local refreshed_at = tonumber(redis.call("get", KEYS[2]))
if refreshed_at == nil then
	refreshed_at = 0
end

-- 时间相隔超过一个间隔，写入修改时间，写入当前token
if now - refreshed_at > per then
	local per_num = math.floor((now - refreshed_at)/per)
	refreshed_at = refreshed_at + per * per_num
	tokens = math.min(cap, tokens + per_num * rate)
end

-- 本次请求token数量是否足够
local allowed = tokens >= requested
-- 桶剩余数量
if allowed then
	tokens = tokens - requested
end

-- 设置剩余token、时间间隔
redis.call("setex", KEYS[1], ttl, tokens)
redis.call("setex", KEYS[2], ttl, refreshed_at)

return allowed
*/

const (
	script = `
local rate = tonumber(ARGV[1])
local per = tonumber(ARGV[2])
local cap = tonumber(ARGV[3])
local ttl = tonumber(ARGV[4])
local now = tonumber(ARGV[5])
local requested = tonumber(ARGV[6])
local tokens = tonumber(redis.call("get", KEYS[1]))
if tokens == nil then tokens = cap end
local refreshed_at = tonumber(redis.call("get", KEYS[2]))
if refreshed_at == nil then refreshed_at = 0 end
if now - refreshed_at > per then
	local per_num = math.floor((now - refreshed_at)/per)
	refreshed_at = refreshed_at + per * per_num
	tokens = math.min(cap, tokens + per_num * rate)
end
local allowed = tokens >= requested
if allowed then tokens = tokens - requested end
redis.call("setex", KEYS[1], ttl, tokens)
redis.call("setex", KEYS[2], ttl, refreshed_at)
return allowed`
)

const (
	redisAliveNo  = 0
	redisAliveYes = 1
)

type TokenLimiter struct {
	rate                   int            // 速率 每per秒生成rate个
	per                    int            // 时间间隔
	burst                  int            // 容量
	ttl                    int            // 超时时间，设为填满桶时间的2倍
	store                  redis.Cmdable  // redisCli
	tokenKey, timestampKey string         // redisKey
	redisAlive             int32          // redis是否存活
	pingInterval           time.Duration  // ping间隔
	localLimiter           *xrate.Limiter // 本地令牌桶
}

// 令牌桶 默认使用Redis用作存储，为分布式令牌桶
// 如Redis异常，则使用本地令牌桶，不会让服务出现单点问题，且继续提供服务
// Redis异常后，会开启Redis存活校验，每隔pingInterval校验一次Redis，默认1s
func NewTokenLimiter(rate, burst int, store redis.Cmdable, key string, opts ...Option) *TokenLimiter {
	cfg := buildConfig(opts)
	if rate < 1 || burst < 1 || burst < rate {
		panic("per、rate、burst有误")
	}

	per := int(cfg.per / time.Second)
	return &TokenLimiter{
		rate:         rate,
		burst:        burst,
		per:          per,
		ttl:          burst * per * 2 / rate, // 填满桶时间的2倍
		store:        store,
		tokenKey:     fmt.Sprintf(cfg.tokenFormat, key),
		timestampKey: fmt.Sprintf(cfg.timestampFormat, key),
		pingInterval: cfg.pingInterval,
		redisAlive:   redisAliveYes,
		localLimiter: xrate.NewLimiter(xrate.Every(cfg.per/time.Duration(rate)), burst),
	}
}

func (lim *TokenLimiter) Allow(ctx context.Context) bool {
	return lim.AllowN(ctx, time.Now(), 1)
}

func (lim *TokenLimiter) AllowN(ctx context.Context, now time.Time, n int) bool {
	return lim.reserveN(ctx, now, n)
}

func (t *TokenLimiter) reserveN(ctx context.Context, now time.Time, n int) bool {
	if atomic.LoadInt32(&t.redisAlive) == redisAliveNo {
		return t.localLimiter.AllowN(now, n)
	}

	resp, err := t.store.Eval(ctx,
		script,
		[]string{t.tokenKey, t.timestampKey},
		cast.ToString(t.rate),
		cast.ToString(t.per),
		cast.ToString(t.burst),
		cast.ToString(t.ttl),
		cast.ToString(now.Unix()),
		cast.ToString(n),
	).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("fail to use rate limiter: %s, use in-process limiter for rescue", err)
		go t.pingRedis()
		return t.localLimiter.AllowN(now, n)
	}

	if errors.Is(err, redis.Nil) {
		return false
	}
	code, _ := resp.(int64)
	return code == 1
}

func (t *TokenLimiter) pingRedis() {
	// 原为存活，改为不存活
	// 仅允许一个pingRedis()启动
	if !atomic.CompareAndSwapInt32(&t.redisAlive, redisAliveYes, redisAliveNo) {
		return
	}

	ticker := time.NewTicker(t.pingInterval)
	defer ticker.Stop()
	for range ticker.C {
		if dbRedis.Ping(t.store) {
			atomic.StoreInt32(&t.redisAlive, redisAliveYes)
			return
		}
	}
}
