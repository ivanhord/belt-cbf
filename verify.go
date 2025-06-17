package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
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

func VerifyMessagesFromFile(path string) error {
	decoder := charmap.Windows1251.NewDecoder()
	encoder := charmap.Windows1251.NewEncoder()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("чтение %s: %w", path, err)
	}

	decoded, err := decoder.Bytes(data)
	if err != nil {
		return fmt.Errorf("декодирование Windows-1251: %w", err)
	}

	messages, err := ParseMessages(string(decoded))
	if err != nil {
		return fmt.Errorf("парсинг сообщений: %w", err)
	}

	found := false
	for _, m := range messages {
		fmt.Printf("метка=%s, сообщение=\"%s\", хэш=%s\n", m.Label, m.Message, m.HashHex)
		msgBytes, err := encoder.Bytes([]byte(m.Message))
		if err != nil {
			return fmt.Errorf("кодировка обратно в 1251: %w", err)
		}
		hash, err := BeltHashGo([]byte(msgBytes))
		if err != nil {
			return fmt.Errorf("хэширование: %w", err)
		}

		expected, err := hex.DecodeString(strings.TrimPrefix(m.HashHex, "0x"))

		if err != nil {
			return fmt.Errorf("декод хэша: %w", err)
		}

		if len(expected) != len(hash) {
			continue
		}
		fmt.Printf("хэш1=%s\n", hash)
		fmt.Printf("хэш2=%s\n", expected)

		matched := true
		for i := range hash {
			print(hash[i])
			print(expected[i])
			if hash[i] != expected[i] {
				matched = false
				break
			}
		}

		if matched {
			fmt.Printf("✅ Сообщение %s корректно:\n%s\n\n", m.Label, m.Message)
			answer := fmt.Sprintf("Сообщение %s корректно:\n%s\n\n", m.Label, m.Message)

			encoded, err := encoder.String(answer)
			if err != nil {
				return fmt.Errorf("ошибка кодирования: %w", err)
			}

			// Записываем в файл
			err = os.WriteFile("message.txt", []byte(encoded), 0644)
			if err != nil {
				return fmt.Errorf("запись в файл: %w", err)
			}
			found = true
			break
		}
	}

	if !found {
		fmt.Println("❌ Ни одно сообщение не совпало с хэшом.")
	}
	return nil
}
