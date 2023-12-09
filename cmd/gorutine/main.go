package main

import (
	"fmt"
	"net/http"
)

func main() {
	var (
		urls = []string{"http://www.googl.com", "https://hoge"}
		done = make(chan interface{})
	)

	defer close(done)

	Do := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		// 並列化する
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				// エラーの処理の責任をここで持たせている
				if err != nil {
					fmt.Println(err)
					continue
				}
				// goroutineをさばく
				for {
					select {
					case responses <- resp:
						fmt.Println("hogehoge1")
					case <-done:
						fmt.Println("hogehoge2")
						return
					}
				}
			}
		}()
		return responses
	}

	for response := range Do(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}
