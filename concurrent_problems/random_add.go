package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func randomAdd() {
	const N = 5

	wg := sync.WaitGroup{}
	wg.Add(N)

	res := make(chan int)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			a := rand.Intn(100)
			b := rand.Intn(100)
			fmt.Println(a + b)
			res <- a + b
		}()
	}

	// 必须开一个协程
	// 否则阻塞在 wg.Wait(), for range res 无法接收数据
	go func() {
		wg.Wait()
		close(res)
	}()

	ans := 0

	for num := range res {
		if ans < num {
			ans = num
		}
	}

	fmt.Println("ans is: ", ans)
}
