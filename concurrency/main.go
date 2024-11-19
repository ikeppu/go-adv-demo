package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	t := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			getHttpCode()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(t))
}

func getHttpCode() {
	resp, err := http.Get("https://google.com")

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		return
	}

	fmt.Printf("StatusCode %d \n", resp.StatusCode)
}
