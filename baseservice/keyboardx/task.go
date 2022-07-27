package keyboardx

import "strings"

type TaskService interface {
	AddFullyTask(key string, f Task) TaskService
	AddPrefixTask(key string, f Task) TaskService
	Match(key string) Task
}

type prefixTask struct {
	key      string
	taskFunc Task
}

type taskService struct {
	fullyTasks  map[string]Task
	prefixTasks []prefixTask
}

func NewTaskService() TaskService {
	return &taskService{
		fullyTasks:  make(map[string]Task),
		prefixTasks: make([]prefixTask, 0),
	}
}

func (t *taskService) AddFullyTask(key string, f Task) TaskService {
	t.fullyTasks[key] = f
	return t
}

func (t *taskService) AddPrefixTask(key string, f Task) TaskService {
	t.prefixTasks = append(t.prefixTasks, prefixTask{
		key:      key,
		taskFunc: f,
	})
	return t
}

func (t *taskService) Match(key string) Task {
	if f, ok := t.fullyTasks[key]; ok {
		return f
	}
	for _, v := range t.prefixTasks {
		if strings.HasPrefix(key, v.key) {
			return v.taskFunc
		}
	}
	return nil
}
