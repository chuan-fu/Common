package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	log "github.com/chuan-fu/Common/zlog"
	"github.com/panjf2000/ants/v2"
)

var (
	pnum  = 10000
	gonum = 1000000
)

func main() {
	usego() // 分配的内存 = 378997KB, GC的次数 = 7
	// useAnts() // 分配的内存 = 9575KB, GC的次数 = 8
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
