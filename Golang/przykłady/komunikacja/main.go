package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	doneCh := make(chan int)
	errCh := make(chan error)

	go produce(ctx, doneCh, errCh)
	consume(doneCh, errCh)
}

func produce(ctx context.Context, c chan<- int, e chan<- error) {
	t := time.NewTicker(time.Second / 3)
	future := time.After(time.Second / 2)
loop:
	for range t.C {
		select {
		case <-ctx.Done():
			e <- errors.New("timeout")
			break loop
		case <-future:
			c <- 1
		default: // !
			c <- 2
		}
	}
	close(c)
	close(e)
}

func consume(c <-chan int, e <-chan error) {
	for {
		select {
		case done := <-c:
			fmt.Println("consumed:", done)
		case e := <-e:
			fmt.Println("job done:", e)
			return
		}
	}
}
