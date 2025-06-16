package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type MessageHash struct {
	Label   string // "А", "Б", "В"
	Message string
	HashHex string
}

func ParseMessages(input string) ([]MessageHash, error) {
	// Разделение по блокам: А., Б., В.
	re := regexp.MustCompile(`(?m)^([А-Я])\.\s+Сообщение:\s*\n?(.+?)\n+Хэш-значение:\s*(0x[0-9A-Fa-f]{64})`)
	matches := re.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return nil, fmt.Errorf("не удалось найти сообщения")
	}

	var result []MessageHash
	for _, m := range matches {
		label := strings.TrimSpace(m[1])
		msg := strings.Trim(m[2], "«»\"“”' \n\t")
		hash := strings.TrimSpace(m[3])
		result = append(result, MessageHash{
			Label:   label,
			Message: msg,
			HashHex: hash,
		})
	}

	return result, nil
}

func VerifyMessagesFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("чтение %s: %w", path, err)
	}

	messages, err := ParseMessages(string(data))
	if err != nil {
		return fmt.Errorf("парсинг сообщений: %w", err)
	}

	found := false
	for _, m := range messages {
		fmt.Printf("метка=%s, сообщение=\"%s\", хэш=%s\n", m.Label, m.Message, m.HashHex)
		hash, err := BeltHashGo([]byte(m.Message))
		if err != nil {
			return fmt.Errorf("хэширование: %w", err)
		}

		expected, err := hex.DecodeString(strings.TrimPrefix(m.HashHex, "0x"))
		fmt.Printf("хэш=%s\n", hash)
		if err != nil {
			return fmt.Errorf("декод хэша: %w", err)
		}

		if len(expected) != len(hash) {
			continue
		}

		matched := true
		for i := range hash {
			if hash[i] != expected[i] {
				matched = false
				break
			}
		}

		if matched {
			fmt.Printf("✅ Сообщение %s корректно:\n%s\n\n", m.Label, m.Message)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("❌ Ни одно сообщение не совпало с хэшом.")
	}
	return nil
}
