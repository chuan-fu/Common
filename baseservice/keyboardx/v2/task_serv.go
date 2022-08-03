package keyboardx

import (
	"fmt"
	"strings"

	"github.com/chuan-fu/Common/util"
)

type TaskService interface {
	TaskKeyList() []string // key列表

	DefaultTask(f Task) TaskService                               // 默认任务
	AddFullyTask(key, desc string, f Task) TaskService            // 全匹配任务
	AddFullyTasks(keys []string, desc string, f Task) TaskService // 全匹配任务
	AddPrefixTask(key, desc string, f Task) TaskService           // 前缀匹配任务
	AddPrefixTasks(key []string, desc string, f Task) TaskService // 前缀匹配任务
	MatchTask(key string) Task                                    // 匹配

	AddHelpTask(key, desc string) TaskService           // 添加功能介绍
	AddHelpTasks(key []string, desc string) TaskService // 添加功能介绍

}

type prefixTask struct {
	key      string
	taskFunc Task
}

type helpShow struct {
	k []string
	v string
}

type taskService struct {
	defaultTask Task
	fullyTasks  map[string]Task
	prefixTasks []prefixTask
	helpShow    []helpShow
}

func NewTaskService() TaskService {
	return &taskService{
		fullyTasks:  make(map[string]Task),
		prefixTasks: make([]prefixTask, 0),
	}
}

func (t *taskService) TaskKeyList() []string {
	list := make([]string, 0, len(t.helpShow))
	for k := range t.helpShow {
		list = append(list, t.helpShow[k].k...)
	}
	return list
}

func (t *taskService) DefaultTask(f Task) TaskService {
	t.defaultTask = f
	return t
}

func (t *taskService) AddFullyTask(key, desc string, f Task) TaskService {
	t.helpShow = append(t.helpShow, helpShow{[]string{key}, desc})
	t.fullyTasks[key] = f
	return t
}

func (t *taskService) AddFullyTasks(keys []string, desc string, f Task) TaskService {
	t.helpShow = append(t.helpShow, helpShow{keys, desc})
	for k := range keys {
		t.fullyTasks[keys[k]] = f
	}
	return t
}

func (t *taskService) AddPrefixTask(key, desc string, f Task) TaskService {
	t.helpShow = append(t.helpShow, helpShow{[]string{key}, desc})
	t.prefixTasks = append(t.prefixTasks, prefixTask{
		key:      key,
		taskFunc: f,
	})
	return t
}

func (t *taskService) AddPrefixTasks(keys []string, desc string, f Task) TaskService {
	t.helpShow = append(t.helpShow, helpShow{keys, desc})
	for k := range keys {
		t.prefixTasks = append(t.prefixTasks, prefixTask{
			key:      keys[k],
			taskFunc: f,
		})
	}
	return t
}

func (t *taskService) MatchTask(key string) Task {
	if f, ok := t.fullyTasks[key]; ok {
		return f
	}
	for k := range t.prefixTasks {
		v := &t.prefixTasks[k]
		if strings.HasPrefix(key, v.key) {
			return v.taskFunc
		}
	}
	return t.defaultTask
}

func (t *taskService) AddHelpTask(key, desc string) TaskService {
	t.AddFullyTask(key, desc, t.help())
	return t
}

func (t *taskService) AddHelpTasks(keys []string, desc string) TaskService {
	t.AddFullyTasks(keys, desc, t.help())
	return t
}

func (t *taskService) help() Task {
	return NewHandleTask(func(s string) (isEnd bool, err error) {
		f := util.NewFmtList()
		for k := range t.helpShow {
			f.Add(fmt.Sprintf("%s -> %s", strings.Join(t.helpShow[k].k, ", "), t.helpShow[k].v))
		}
		fmt.Println(f.String())
		return
	})
}
