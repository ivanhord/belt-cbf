package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/decrypt", decryptHandler)

	log.Println("🚀 Decrypt service listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
