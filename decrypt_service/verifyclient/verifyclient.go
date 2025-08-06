package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestToVerify struct {
	Message string `json:"message"`
}

type ResponseOnVerify struct {
	Status bool   `json:"status"`
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

func CallVerifyService(message string) (*ResponseOnVerify, error) {
	// Формируем JSON-тело
	reqBody := RequestToVerify{
		Message: message,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации: %w", err)
	}

	// Отправляем POST-запрос на verify_service
	resp, err := http.Post("http://verify_service:8881/verify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к verify_service: %w", err)
	}
	defer resp.Body.Close()

	// Декодируем ответ
	var result ResponseOnVerify
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &result, nil
}
