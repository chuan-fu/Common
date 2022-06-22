package cdao

/* hmsetScript hash写入
redis.call("hmset", KEYS[1], unpack(ARGV))
if KEYS[2] then
	redis.call("expire", KEYS[1], KEYS[2])
end
*/

/* zaddScript 覆盖式写入
if redis.call("exists", KEYS[1]) == 1 then
	redis.call("del", KEYS[1])
end
redis.call("zadd", KEYS[1], unpack(ARGV))
if KEYS[2] then
	redis.call("expire", KEYS[1], KEYS[2])
end
return 1
*/

/* zgetallScript 获取所有
local num = redis.call("zcard", KEYS[1])
if num > 0 then
	return redis.call("zrange", KEYS[1], 0, num)
end
*/

/* saddScript 覆盖式写入
if redis.call("exists", KEYS[1]) == 1 then
	redis.call("del", KEYS[1])
end
redis.call("sadd", KEYS[1], unpack(ARGV))
if KEYS[2] then
	redis.call("expire", KEYS[1], KEYS[2])
end
return 1
*/

/* setbits 写入多个bit
for i, v in ipairs(ARGV) do
	redis.call("setbit", KEYS[1], v, 1)
end
if KEYS[2] then
	redis.call("expire", KEYS[1], KEYS[2])
end
return 1
*/

/* getbits 写入多个bit
if redis.call("exists", KEYS[1]) == 0 then
	return
end
local resp = {}
for i, v in ipairs(ARGV) do
	if v >= 0 then
		resp[i] = redis.call("getbit", KEYS[1], v)
	else
		resp[i] = 0
	end
end
return resp
*/

const (
	hmsetScript   = `redis.call("hmset", KEYS[1], unpack(ARGV)) if KEYS[2] then redis.call("expire", KEYS[1], KEYS[2]) end`
	zaddScript    = `if redis.call("exists", KEYS[1]) == 1 then redis.call("del", KEYS[1]) end redis.call("zadd", KEYS[1], unpack(ARGV)) if KEYS[2] then redis.call("expire", KEYS[1], KEYS[2]) end return 1`
	saddScript    = `if redis.call("exists", KEYS[1]) == 1 then redis.call("del", KEYS[1]) end redis.call("sadd", KEYS[1], unpack(ARGV)) if KEYS[2] then redis.call("expire", KEYS[1], KEYS[2]) end return 1`
	zgetallScript = `local num = redis.call("zcard", KEYS[1]) if num > 0 then return redis.call("zrange", KEYS[1], 0, num) end`

	setBitsScript = `for i, v in ipairs(ARGV) do redis.call("setbit", KEYS[1], v, 1) end if KEYS[2] then redis.call("expire", KEYS[1], KEYS[2]) end return 1`
	getBitsScript = `if redis.call("exists", KEYS[1]) == 0 then return end local resp = {} for i, v in ipairs(ARGV) do resp[i] = redis.call("getbit", KEYS[1], v) end return resp`
)
