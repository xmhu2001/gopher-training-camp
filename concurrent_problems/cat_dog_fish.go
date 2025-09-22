package main

import (
	"fmt"
	"sync"
	"time"
)

func cat_dog_fish() {

	wg := sync.WaitGroup{}
	wg.Add(3)

	catChannel := make(chan struct{})
	dogChannel := make(chan struct{})
	fishChannel := make(chan struct{})

	go func() {
		defer wg.Done()
		for {
			<-catChannel
			fmt.Println("cat")
			time.Sleep(time.Second)
			dogChannel <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			<-dogChannel
			fmt.Println("dog")
			time.Sleep(time.Second)
			fishChannel <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			<-fishChannel
			fmt.Println("fish")
			time.Sleep(time.Second)
			catChannel <- struct{}{}
		}
	}()
	catChannel <- struct{}{}
	wg.Wait()
}
