package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var wg sync.WaitGroup

func food(food chan int) {
	food <- rand.Intn(100)
	wg.Done()
}

func selectFunc() {

}

func main() {
	test := make(chan int)
	quit := make(chan int)
	wg.Add(3)
	go food(test)
	go food(test)
	go food(test)

	for i := 0; i < 3; i++ {
		select {
		case val := <-test:
			fmt.Println("Got some value from channel:", val)
		case <-quit:
			return
		}
	}
	wg.Wait()
	close(test)

}
