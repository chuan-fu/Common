package util

import (
	"context"
	"log"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestCtx(t *testing.T) {
	ctx := context.Background()

	ctx2, cancel2 := context.WithCancel(ctx)
	go func() {
		<-ctx2.Done()
		log.Println("close2")
	}()

	ctx3, cancel3 := context.WithCancel(ctx2)
	go func() {
		<-ctx3.Done()
		log.Println("close3")
	}()

	ctx4, cancel4 := context.WithCancel(ctx3)
	go func() {
		<-ctx4.Done()
		log.Println("close4")
	}()

	time.Sleep(2 * time.Second)
	log.Println("cancel3")
	cancel3()

	time.Sleep(2 * time.Second)
	log.Println("cancel2")
	cancel2()

	time.Sleep(2 * time.Second)
	log.Println("cancel4")
	cancel4()

	time.Sleep(2 * time.Second)
}

func TestErrGroup(t *testing.T) {
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return nil
	})
}
