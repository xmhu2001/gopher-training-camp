package main

import (
	"fmt"
	"sync"
)

const (
	totalJobs int = 1000
	producers int = 10
	consumers int = 5
)

func producerConsumer() {
	queue := make(chan int, totalJobs)

	wg := sync.WaitGroup{}
	wg.Add(producers)

	for i := 0; i < producers; i++ {
		go func() {
			defer wg.Done()
			for j := 1; j <= totalJobs; j++ {
				if j%producers == i {
					queue <- j
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	cg := sync.WaitGroup{}
	cg.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func() {
			defer cg.Done()
			for v := range queue {
				fmt.Printf("consumer %d handled %d\n", i, v)
			}
		}()
	}

	cg.Wait()
}
