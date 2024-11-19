package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func main() {
	http.HandleFunc("/hello", hello)

	fmt.Println("Server is listening on port 8081")
	// Default server mux
	http.ListenAndServe(":8081", nil)
}
