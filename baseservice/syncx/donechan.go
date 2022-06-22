package syncx

import "sync"

type DoneChan interface {
	Wait()
	Done() chan struct{}
	Close()
}

type doneChan struct {
	sync.Once
	done chan struct{}
}

func NewDoneChan() DoneChan {
	return &doneChan{
		done: make(chan struct{}),
	}
}

func (d *doneChan) Wait() {
	<-d.done
}

func (d *doneChan) Done() chan struct{} {
	return d.done
}

func (d *doneChan) Close() {
	d.Do(func() {
		close(d.done)
	})
}
