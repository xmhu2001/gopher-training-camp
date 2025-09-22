package main

import (
	"fmt"
	"sync"
)

func sumOddEven() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	oddChannel := make(chan int)
	evenChannel := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		sum := 0
		for num := range oddChannel {
			sum += num
		}
		fmt.Println("oddSum:", sum)
	}()

	go func() {
		defer wg.Done()
		sum := 0
		for num := range evenChannel {
			sum += num
		}
		fmt.Println("evenSum:", sum)
	}()

	for _, v := range nums {
		if v%2 == 1 {
			oddChannel <- v
		} else {
			evenChannel <- v
		}
	}

	close(oddChannel)
	close(evenChannel)
	wg.Wait()
}
