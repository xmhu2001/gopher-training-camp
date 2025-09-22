package main

import (
	"fmt"
	"sync"
)

const (
	n     int = 5
	limit int = 50
)

func nGoroutine() {
	channels := make([]chan int, n)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		channels[i] = make(chan int)
	}

	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			next := (idx + 1) % n
			for j := idx + 1; j <= limit; j += n {
				<-channels[idx]
				fmt.Println(j)
				if j == limit {
					continue
				}
				channels[next] <- 1
			}
		}(i)
	}

	channels[0] <- 1
	wg.Wait()
}
