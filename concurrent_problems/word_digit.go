package main

import (
	"fmt"
	"sync"
	"time"
)

func word_digit() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	wordChannel := make(chan struct{})
	digitChannel := make(chan struct{})

	go func() {
		defer wg.Done()
		for i := 0; i < 26; i++ {
			<-digitChannel
			fmt.Println(i)
			time.Sleep(time.Millisecond * 500)
			wordChannel <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 26; i++ {
			<-wordChannel
			fmt.Printf("%c\n", rune(i+'a'))
			time.Sleep(time.Millisecond * 500)
			digitChannel <- struct{}{}
		}
	}()

	digitChannel <- struct{}{}
	wg.Wait()
}
