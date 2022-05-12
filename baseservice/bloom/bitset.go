package bloom

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

/* setScript
for _, offset in ipairs(ARGV) do
	redis.call("setbit", KEYS[1], offset, 1)
end
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
	setScript   = `for _, offset in ipairs(ARGV) do redis.call("setbit", KEYS[1], offset, 1) end`
	checkScript = `for _, offset in ipairs(ARGV) do if tonumber(redis.call("getbit", KEYS[1], offset)) == 0 then return 0 end end return 1`

	isExist = 1 // 存在
)

var (
	ErrTooLargeOffset = errors.New("too large offset")
	ErrAssertion      = errors.New("assertion err")
)

type BitSetProvider interface {
	check([]uint) (bool, error)
	set([]uint) error
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
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	if exists, ok := result.(int64); ok {
		return exists == isExist, nil
	} else {
		return false, ErrAssertion
	}
}

func (r *redisBitSet) set(offsets []uint) error {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return err
	}
	_, err = r.rCli.Eval(context.Background(), setScript, []string{r.key}, args).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
