package main

import (
	"fmt"
	"sync"
)

const num int = 100

func printById() {

	wg := sync.WaitGroup{}

	chs := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
	}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			next := (id + 1) % 10
			for j := id; j <= num; j += 10 {
				<-chs[id]
				fmt.Printf("process %d: %d\n", id, j)
				if j == num {
					continue
				}
				chs[next] <- 1
			}
		}(i)
	}

	chs[0] <- 1
	wg.Wait()
}
