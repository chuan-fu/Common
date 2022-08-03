package keyboardx

type HandleFunc func(s string) (isEnd bool, err error)

type HandleTask struct {
	BaseTask
	handleFunc HandleFunc
}

func (h *HandleTask) Handle(s string) (err error) {
	if h.handleFunc != nil {
		h.IsEndField, err = h.handleFunc(s)
	}
	return
}

func NewHandleTask(f HandleFunc) Task {
	return &HandleTask{handleFunc: f}
}

func NewExitTask() Task {
	return &HandleTask{
		handleFunc: func(s string) (isEnd bool, err error) {
			return true, nil
		},
	}
}
