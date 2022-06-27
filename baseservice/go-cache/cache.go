package cache

import (
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var ExistErr = errors.New("key already exists")

type Item struct {
	Object     interface{}
	Expiration int64
}

// 是否过期
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

func ToItem(k interface{}) *Item {
	if v, ok := k.(*Item); ok {
		return v
	}
	return new(Item)
}

const (
	// 永不过期
	NoExpiration time.Duration = -1
	// 使用默认过期时间
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration time.Duration // 默认过期时间
	items             sync.Map
	onEvicted         func(string, interface{})
	janitor           *janitor
}

// 写入
// d > 0 使用d
// d == 0 使用 defaultExpiration
// d < 0 永不过期
func (c *cache) Set(k string, x interface{}, d time.Duration) {
	c.set(k, x, d)
}

// 使用默认过期的写入
func (c *cache) SetDefault(k string, x interface{}) {
	c.set(k, x, DefaultExpiration)
}

// 使用默认过期的写入
func (c *cache) SetForever(k string, x interface{}) {
	c.set(k, x, NoExpiration)
}

func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items.Store(k, &Item{
		Object:     x,
		Expiration: e,
	})
}

// 不存在则添加 返回是否添加成功
func (c *cache) Add(k string, x interface{}, d time.Duration) bool {
	if c.exists(k) {
		return false
	}
	c.set(k, x, d)
	return true
}

func (c *cache) Get(k string) (interface{}, bool) {
	itemInter, found := c.items.Load(k)
	if !found {
		return nil, false
	}
	item := ToItem(itemInter)
	if item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// 获取数据、超时时间
func (c *cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	itemInter, found := c.items.Load(k)
	if !found {
		return nil, time.Time{}, false
	}
	item := ToItem(itemInter)
	if item.Expired() {
		return nil, time.Time{}, true
	}
	if item.Expiration == 0 {
		return item.Object, time.Time{}, true
	}
	return item.Object, time.Unix(0, item.Expiration), true
}

func (c *cache) exists(k string) bool {
	itemInter, found := c.items.Load(k)
	if !found {
		return false
	}

	item := ToItem(itemInter)
	if item.Expired() {
		return false
	}
	return true
}

// 删除缓存
func (c *cache) Delete(k string) {
	if c.onEvicted != nil {
		if vInter, found := c.items.Load(k); found {
			c.items.Delete(k)
			c.onEvicted(k, ToItem(vInter).Object)
			return
		}
	}
	c.items.Delete(k)
}

type kav struct {
	key   string
	value interface{}
}

func (c *cache) DeleteExpired() {
	now := time.Now().UnixNano()
	delList := make([]*kav, 0)
	c.items.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v := ToItem(value)
		if v.Expiration > 0 {
			if now > v.Expiration {
				delList = append(delList, &kav{key: k, value: v.Object})
			}
		}
		return true
	})
	for k := range delList {
		v := delList[k]
		c.items.Delete(v.key)
		if c.onEvicted != nil {
			c.onEvicted(v.key, v.value)
		}
	}
}

// 删除时触发函数
func (c *cache) OnEvicted(f func(string, interface{})) {
	c.onEvicted = f
}

// 缓存列表
func (c *cache) Items() map[string]Item {
	m := make(map[string]Item, 0)
	now := time.Now().UnixNano()
	c.items.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v := ToItem(value)
		if v.Expiration > 0 {
			if now > v.Expiration {
				return true
			}
		}
		m[k] = *v
		return true
	})
	return m
}

func (c *cache) Flush() {
	c.items = sync.Map{}
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             sync.Map{},
	}
	for k := range m {
		v := m[k]
		c.items.Store(k, &v)
	}
	return c
}

func newCacheWithJanitor(de, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, make(map[string]Item))
}

func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}
