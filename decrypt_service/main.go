package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/decrypt", decryptHandler)

	log.Println("🚀 Decrypt service listening on :8880")
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Fatal(err)
	}
}
