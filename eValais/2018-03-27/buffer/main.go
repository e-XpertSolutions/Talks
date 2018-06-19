package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doSomething(buf []int) {
	if len(buf) == 0 {
		return
	}
	for _, v := range buf {
		fmt.Printf("%d ", v)
	}
	fmt.Print("\n")
}

func routine(wg *sync.WaitGroup, ch <-chan int, stopc <-chan bool) {
	defer wg.Done()

	var buf []int
	for {
		select {
		case i := <-ch:
			buf = append(buf, i)
			if len(buf) > 10 {
				doSomething(buf)
				buf = []int{}
			}
		case <-time.After(500 * time.Millisecond):
			doSomething(buf)
			buf = []int{}
		case <-stopc:
			return
		}
	}
}

func main() {
	ch := make(chan int, 2)
	stopc := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)

	go routine(&wg, ch, stopc)

	for i := 0; i < 10; i++ {
		ch <- rand.Intn(100)
		time.Sleep(time.Duration(rand.Intn(700)) * time.Millisecond)
	}

	stopc <- true
	wg.Wait()

	close(ch)
	close(stopc)
}
