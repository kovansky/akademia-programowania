package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const count = 5

type storage struct {
	sync.Mutex

	i       int
	mutexI  int
	atomicI atomic.Int64
}

func (s *storage) add(i int) {
	s.i += i

	s.atomicI.Add(int64(i))

	s.Lock()
	s.mutexI += i
	s.Unlock()
}

func main() {
	s := storage{}
	for i := 0; i < count*1000; i++ {
		go s.add(1)
	}
	time.Sleep(time.Second)
	fmt.Println("I:", s.i)
	fmt.Println("Atomic I", s.atomicI.Load())
	fmt.Println("Mutex I", s.mutexI)

	for i := 0; i < count; i++ {
		go func() {
			fmt.Println("func :", i) // loop variable i captured by func literal
		}()
	}
	time.Sleep(time.Second / 2)

	for i := 0; i < count; i++ {
		go func(i int) { // var shadowing = copy
			fmt.Println("BRR:", i+1)
		}(i)
		go brr(i) // passing by value = copy
	}
	time.Sleep(time.Second / 2)

	intCh := make(chan int, 10)
	defer close(intCh)

	go foo(intCh)

	time.Sleep(time.Second / 2)
	for v := range intCh {
		fmt.Println("got", v)
	}

	x, isOpen := <-intCh
	fmt.Println(x, isOpen)
	intCh <- 1337
}

func brr(i int) {
	fmt.Println("brr:", i)
}

func foo(c chan<- int) {
	for i := 0; i < count; i++ {
		c <- i
	}
	close(c)
}
