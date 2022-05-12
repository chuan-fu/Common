package redislock

import (
	"time"
)

/* delLockScript
local v2 = redis.call("get", KEYS[1])
if v2 then
	if v2 == ARGV[1] then
		redis.call("del", KEYS[1])
		return 1
	end
	return -1
end
return 0
*/

/* extendLockScript
local v2 = redis.call("get", KEYS[1])
if v2 then
	if v2 == ARGV[1] then
		redis.call("pexpire", KEYS[1], ARGV[2])
		return 1
	end
	return -1
else
	redis.call("set", KEYS[1], ARGV[1], "px", ARGV[2])
	return 1
end
*/

const (
	extendLockScript     = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("pexpire", KEYS[1], ARGV[2]) return 1 end return -1 else redis.call("set", KEYS[1], ARGV[1], "px", ARGV[2]) return 1 end`
	delLockScript        = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("del", KEYS[1]) return 1 end return -1 end return 0`
	DelLockStatusNotOwn  = -1 // 非本人的锁【锁已过期，且已被抢占】
	DelLockStatusExpired = 0  // 锁已过期
	DelLockStatusSuccess = 1  // 删除成功
)

func formatMs(dur time.Duration) int64 {
	if dur <= 0 {
		return -1
	}
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}
