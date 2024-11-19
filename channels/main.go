package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	// Create channel
	code := make(chan int)
	//
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {

		wg.Add(1)
		go func() {
			getHttpCode(code)
			wg.Done()
		}()
	}

	// Read from channel multiple
	go func() {
		wg.Wait()
		close(code)
	}()
	for res := range code {
		fmt.Printf("StatusCode Main %d \n", res)
	}

	// Read from channel once
	// res := <-code

}

func getHttpCode(codeCh chan int) {
	resp, err := http.Get("https://google.com")

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}

	// Write in channel
	codeCh <- resp.StatusCode
}
