package cron

import (
	"github.com/chuan-fu/Common/baseservice/mutex"
	"github.com/chuan-fu/Common/zlog"
	"github.com/robfig/cron"
)

func defaultCondition() bool { return true }

func RunCronTask(tasks ...*cronTask) {
	if len(tasks) == 0 {
		return
	}

	var hasTask bool
	crontab := cron.New()

	for k := range tasks {
		v := tasks[k]
		if v == nil || !v.condition() {
			continue
		}

		hasTask = true
		log.Infof("RunCronTask AddFunc spce:%s , name:%s", v.spec, v.name)
		err := crontab.AddFunc(v.spec, func() func() {
			if v.onceKey != "" {
				distributedOnce := mutex.NewDistributedOnce(v.onceKey, v.opts...)
				return func() {
					log.Infof("RunCronTask TaskRun once spce:%s , name:%s", v.spec, v.name)
					distributedOnce.Do(v.cmd)
				}
			}
			return func() {
				log.Infof("RunCronTask TaskRun spce:%s , name:%s", v.spec, v.name)
				v.cmd()
			}
		}())
		if err != nil {
			log.Fatal(err)
		}
	}

	if hasTask {
		log.Info("RunCronTask Start")
		crontab.Start()
	}
}

type cronTask struct {
	name string // 任务名称

	spec string
	cmd  func()

	condition func() bool

	onceKey string
	opts    []mutex.Option
}

type Option func(*cronTask)

func WithCondition(condition func() bool) Option {
	return func(once *cronTask) {
		once.condition = condition
	}
}

func WithName(name string) Option {
	return func(once *cronTask) {
		once.name = name
	}
}

func WithMutex(key string, opts ...mutex.Option) Option {
	return func(once *cronTask) {
		once.onceKey = key
		once.opts = opts
	}
}

func NewCronTask(spec string, cmd func(), opts ...Option) *cronTask {
	c := &cronTask{
		spec:      spec,
		cmd:       cmd,
		condition: defaultCondition,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
