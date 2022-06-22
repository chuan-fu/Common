package syncx

import (
	"sync"
	"testing"
)

func TestDoneChanClose(t *testing.T) {
	d := NewDoneChan()

	for i := 0; i < 5; i++ {
		d.Close()
	}
}

func TestDoneChanDone(t *testing.T) {
	var waitGroup sync.WaitGroup
	d := NewDoneChan()

	waitGroup.Add(1)
	go func() {
		<-d.Done()
		waitGroup.Done()
	}()

	for i := 0; i < 5; i++ {
		d.Close()
	}

	waitGroup.Wait()
}
