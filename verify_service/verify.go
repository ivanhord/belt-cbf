// verify.go
package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/ivanhord/belt-cbf/shared/bee2"
)

type MessageHash struct {
	Label   string
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

		// Удаляем лишнее в конце (кавычки, пробелы, управляющие символы)
		msg := strings.TrimRightFunc(m[2], func(r rune) bool {
			switch r {
			case '»', '"', '\r', '\n', '\t', ' ', '“', '”', '«':
				return true
			default:
				return false
			}
		})
		// Удаляем открывающие кавычки и пробелы в начале
		msg = strings.TrimLeft(msg, "«\"“”' \n\t\r")

		hash := strings.TrimSpace(m[3])
		result = append(result, MessageHash{
			Label:   label,
			Message: msg,
			HashHex: hash,
		})
	}

	return result, nil
}

func VerifyMessages(message string) (string, error) {
	//decoder := charmap.Windows1251.NewDecoder()
	encoder := charmap.Windows1251.NewEncoder()
	messages, err := ParseMessages(message)
	if err != nil {
		return "", fmt.Errorf("парсинг сообщений: %w", err)
	}

	for _, m := range messages {
		msgBytes, err := encoder.Bytes([]byte(m.Message))
		if err != nil {
			return "", fmt.Errorf("кодировка обратно в 1251: %w", err)
		}
		hash, err := bee2.BeltHashGo(msgBytes)
		if err != nil {
			return "", fmt.Errorf("хэширование: %w", err)
		}

		expected, err := hex.DecodeString(strings.TrimPrefix(m.HashHex, "0x"))
		if err != nil {
			return "", fmt.Errorf("декод хэша: %w", err)
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
			answer := fmt.Sprintf("Сообщение %s корректно:\n%s", m.Label, m.Message)
			return answer, nil
		}
	}

	return "❌ Ни одно сообщение не совпало с хэшем.", nil
}
