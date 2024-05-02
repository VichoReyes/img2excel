package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server on http://localhost:8080/")
	http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
}
