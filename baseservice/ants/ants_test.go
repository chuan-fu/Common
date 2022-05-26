package ants

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/zlog"
	"github.com/panjf2000/ants/v2"
)

var (
	pnum  = 10000
	gonum = 1000000
)

func TestAnts(t *testing.T) {
	useGlobalAnts()
	// checkAntsClose()
	// usego() // 分配的内存 = 378997KB, GC的次数 = 7
	// useAnts() // 分配的内存 = 9575KB, GC的次数 = 8
}

func useGlobalAnts() {
	NewGlobalPool(20)
	for i := 0; i < 30; i++ {
		go func(index int) {
			GoVoid(func() (err error) {
				time.Sleep(time.Second)
				return
			})
		}(i)
	}
	time.Sleep(5 * time.Second)
}

func checkAntsClose() {
	// ants.WithMaxBlockingTasks(-10)
	p, _ := ants.NewPool(10, ants.WithNonblocking(true), ants.WithPanicHandler(util.DeferFuncLog))
	fmt.Println()
	fmt.Println()

	for i := 0; i < 50; i++ {
		go func(index int) {
			err := p.Submit(func() {
				time.Sleep(time.Second * 5)
				fmt.Println(index)
			})
			if err != nil {
				log.Error(err)
			}
		}(i)
	}
	fmt.Println("over")

	go func() {
		for {
			fmt.Println("waiting =>", p.Waiting())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	t1 := time.Now()
	time.Sleep(3 * time.Minute)
	// p.Release()
	log.Info(time.Now().Sub(t1))
	time.Sleep(10 * time.Second)
	fmt.Println()
	fmt.Println(p.Submit(func() {
		fmt.Println("测试继续写入")
	}))
}

func usego() {
	t1 := time.Now()
	var wg sync.WaitGroup
	wg.Add(gonum)
	for i := 0; i < gonum; i++ {
		go func(index int) {
			// fmt.Println(index)
			time.Sleep(time.Second)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(t1))
	printMemStats("use go")
}

func useAnts() {
	p, _ := ants.NewPool(pnum)
	defer p.Release()

	var wg sync.WaitGroup

	t1 := time.Now()
	for i := 0; i < gonum; i++ {
		wg.Add(1)
		err := p.Submit(func() {
			func(index int) {
				fmt.Println(index)
				time.Sleep(time.Second)
				wg.Done()
			}(i)
		})
		if err != nil {
			log.Error(err)
		}
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(t1))
	printMemStats("use ants")
}

func printMemStats(mag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v：分配的内存 = %vKB, GC的次数 = %v\n", mag, m.Alloc/1024, m.NumGC)
}
