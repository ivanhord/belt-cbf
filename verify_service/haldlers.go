// handlers.go
package main

import (
	"encoding/json"
	"net/http"
)

type RequestToVerify struct {
	Message string `json:"message"`
}

type ResponseOnVerify struct {
	Status bool   `json: "status"`
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is avaible", http.StatusMethodNotAllowed)
		return
	}
	var verify_req RequestToVerify
	if err := json.NewDecoder(r.Body).Decode(&verify_req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	answer, err := VerifyMessages(verify_req.Message)
	if err != nil {
		ver_resp := ResponseOnVerify{
			Status: false,
			Error:  err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ver_resp)
		return
	}
	ver_resp := ResponseOnVerify{
		Status: true,
		Answer: answer,
	}
	json.NewEncoder(w).Encode(ver_resp)
}
