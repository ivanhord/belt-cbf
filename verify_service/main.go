package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/verify", verifyHandler)
	log.Println("🚀 Decrypt service listening on :8881")
	err := http.ListenAndServe(":8881", nil)
	if err != nil {
		log.Fatal(err)
	}

}
