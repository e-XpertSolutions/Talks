package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/satori/go.uuid.v1"
)

type Pool struct {
	uidc chan string // channel of UID
}

func CreatePool(size int) *Pool {
	p := Pool{
		uidc: make(chan string, size),
	}
	for i := 0; i < 4; i++ {
		go p.fillRoutine()
	}
	return &p
}

func (p *Pool) Chan() <-chan string {
	return p.uidc
}

func (p *Pool) Close() {
	close(p.uidc)
}

func (p *Pool) fillRoutine() {
	for {
		p.uidc <- uuid.NewV4().String()
	}
}

var (
	n      = flag.Int("pool.size", 10000, "maximum size of the pool")
	nopool = flag.Bool("pool.disabled", false, "disable pool")
)

const TotalUID = 1000000

func main() {
	flag.Parse()

	if *nopool {
		tic := time.Now()
		for i := 0; i < TotalUID; i++ {
			_ = uuid.NewV4().String()
		}
		fmt.Println("elapsed time: ", time.Since(tic))
		os.Exit(0)
	}

	p := CreatePool(*n)
	defer p.Close()

	time.Sleep(1 * time.Second)

	tic := time.Now()
	for i := 0; i < TotalUID; i++ {
		<-p.Chan()
	}
	fmt.Println("elapsed time: ", time.Since(tic))
}
