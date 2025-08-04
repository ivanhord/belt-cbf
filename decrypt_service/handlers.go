package main

import (
	"encoding/json"
	"net/http"
)

type RequestToDecrypt struct {
	Hex string `json:"hex"`
}

type ResponseOnDecrypt struct {
	Plaintext string `json:"plaintext"`
	Error     string `json:"error,omitempty"`
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is avaible", http.StatusMethodNotAllowed)
		return
	}
	var decr_req RequestToDecrypt
	if err := json.NewDecoder(r.Body).Decode(&decr_req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	plaintext, err := DecryptHex(decr_req.Hex)
	if err != nil {
		decr_resp := ResponseOnDecrypt{Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(decr_resp)
		return
	}
	decr_resp := ResponseOnDecrypt{Plaintext: string(plaintext)}
	json.NewEncoder(w).Encode(decr_resp)

}
