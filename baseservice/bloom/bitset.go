package bloom

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

/* setScript
local state = 0
for _, offset in ipairs(ARGV) do
	if redis.call("setbit", KEYS[1], offset, 1) == 0 then
		state = 1
	end
end
return state
*/

/* checkScript
for _, offset in ipairs(ARGV) do
	if tonumber(redis.call("getbit", KEYS[1], offset)) == 0 then
		return 0
	end
end
return 1
*/

const (
	// setScript   = `for _, offset in ipairs(ARGV) do redis.call("setbit", KEYS[1], offset, 1) end return 0`
	setScript   = `local state = 0 for _, offset in ipairs(ARGV) do if redis.call("setbit", KEYS[1], offset, 1) == 0 then state = 1	end end return state`
	checkScript = `for _, offset in ipairs(ARGV) do if tonumber(redis.call("getbit", KEYS[1], offset)) == 0 then return 0 end end return 1`

	isExist = 1 // 存在
	isAdd   = 1 // 成功添加
)

var (
	ErrTooLargeOffset = errors.New("too large offset")
	ErrAssertion      = errors.New("assertion err")
)

type BitSetProvider interface {
	check([]uint) (bool, error)
	set([]uint) (bool, error)
}

type redisBitSet struct {
	rCli redis.Cmdable
	key  string
	bits uint
}

func newRedisBitSet(rCli redis.Cmdable, key string, bits uint) *redisBitSet {
	return &redisBitSet{
		rCli: rCli,
		key:  key,
		bits: bits,
	}
}

func (r *redisBitSet) buildOffsetArgs(offsets []uint) ([]string, error) {
	args := make([]string, len(offsets))
	for k := range offsets {
		if offsets[k] >= r.bits {
			return nil, ErrTooLargeOffset
		}
		args[k] = fmt.Sprintf("%d", offsets[k])
	}
	return args, nil
}

func (r *redisBitSet) check(offsets []uint) (bool, error) {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return false, err
	}

	result, err := r.rCli.Eval(context.Background(), checkScript, []string{r.key}, args).Result()
	if err != nil {
		return false, err
	}

	if exists, ok := result.(int64); ok {
		return exists == isExist, nil
	} else {
		return false, ErrAssertion
	}
}

func (r *redisBitSet) set(offsets []uint) (bool, error) {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return false, err
	}
	result, err := r.rCli.Eval(context.Background(), setScript, []string{r.key}, args).Result()
	if err != nil {
		return false, err
	}

	if state, ok := result.(int64); ok {
		return state == isAdd, nil
	} else {
		return false, ErrAssertion
	}
}
