package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	// path := flag.String("file", "url.txt", "path to url")
	flag.Parse()

	// file, err := os.ReadFile("C:/Work/go_courses/go-adv-demo/conc_home_work/" + *path)

	// if err != nil {
	// 	panic(err.Error())
	// }

	// urlSlice := strings.Split(string(file), "\n")
	urlSlice := []string{"https://google.com", "https://purpleschool.ru", "htts://ya.ru"}

	resCh := make(chan int)
	errCh := make(chan error)

	for _, url := range urlSlice {
		go ping(url, resCh, errCh)
	}

	for range urlSlice {
		select {
		case errRes := <-errCh:
			fmt.Println(errRes)
		case res := <-resCh:
			fmt.Println(res)

		}

	}
}

func ping(url string, respCh chan int, errCh chan error) {
	fmt.Println(url)
	res, err := http.Get(url)

	if err != nil {
		errCh <- err
		return
	}

	respCh <- res.StatusCode
}
