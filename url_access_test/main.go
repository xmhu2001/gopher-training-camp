package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(url string) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func main() {

	now := time.Now()
	defer func() {
		fmt.Printf("程序耗时： %v\n", time.Since(now))
	}()

	file, _ := os.Open("data.csv")
	defer file.Close()
	goodFile, _ := os.Create("good.csv")
	defer goodFile.Close()
	badFile, _ := os.Create("bad.csv")
	defer badFile.Close()

	reader := csv.NewReader(file)
	badCsv := csv.NewWriter(badFile)
	defer badCsv.Flush()
	goodCsv := csv.NewWriter(goodFile)
	defer goodCsv.Flush()

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// 写入表头
	badCsv.Write(records[0])
	goodCsv.Write(records[0])

	// worker pool
	w := Worker{
		inCh:   make(chan []string),
		goodCh: make(chan []string),
		badCh:  make(chan []string),
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		done: make(chan struct{}),
	}

	go w.Produce(records[1:])

	w.Start(100)

	go w.WriteToFile(badCsv, goodCsv, len(records))

	w.wg.Wait()
	close(w.goodCh)
	close(w.badCh)
	<-w.done
}
