package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CertifiSafe")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Printf("Starting server on :8000...")
	http.ListenAndServe(":8000", nil)
}
