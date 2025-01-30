package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Worker struct {
	inCh   chan []string // 输入通道
	goodCh chan []string // 输出通道
	badCh  chan []string // 输出通道
	client http.Client
	wg     sync.WaitGroup
	done   chan struct{}
}

func (w *Worker) Produce(rows [][]string) {
	for _, row := range rows {
		w.inCh <- row
	}
	close(w.inCh)
}

func (w *Worker) Start(n int) {
	for range n {
		w.wg.Add(1)
		go w.process()
	}
}

func (w *Worker) process() {
	defer w.wg.Done()

	for row := range w.inCh {
		if len(row) < 5 {
			continue
		}

		if w.isValidUrl(row[4]) {
			w.goodCh <- row
		} else {
			w.badCh <- row
		}
	}
}

func (w *Worker) isValidUrl(url_ string) bool {
	url_ = strings.TrimSpace(url_)

	_, err := url.ParseRequestURI(url_)
	if err != nil {
		return false
	}

	resp, err := w.client.Get(url_)

	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}

	return false
}

func (w *Worker) WriteToFile(badCsv, goodCsv *csv.Writer, n int) {
	var finished int = 1 // 初始值是1, 原因为表头没算入

	// 关闭done通道，表明，WriteToFile协程结束
	defer close(w.done)

	for {
		select {
		case row, ok := <-w.goodCh:
			if !ok {
				return
			}

			goodCsv.Write(row)

			finished++

			if finished%10 == 0 {
				fmt.Printf("已完成: %.3f\n", float64(finished)/float64(n))
			}

		case row, ok := <-w.badCh:
			if !ok {
				return
			}

			badCsv.Write(row)

			finished++
			if finished%10 == 0 {
				fmt.Printf("已完成: %.3f\n", float64(finished)/float64(n))
			}
		}
	}
}
