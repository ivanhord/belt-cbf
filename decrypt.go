package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func DecryptFile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("чтение %s: %w", inputPath, err)
	}

	hexCipher := strings.TrimSpace(string(data))
	prefixes := []string{"0x", "0X", "\\x"}
	for _, p := range prefixes {
		if strings.HasPrefix(hexCipher, p) {
			hexCipher = strings.TrimPrefix(hexCipher, p)
			break
		}
	}

	ciphertext, err := hex.DecodeString(hexCipher)
	if err != nil {
		return fmt.Errorf("ошибка декодирования HEX: %w", err)
	}

	key := []byte{
		0x34, 0x87, 0x24, 0xA4, 0xC1, 0xA6, 0x76, 0x67,
		0x15, 0x3D, 0xDE, 0x59, 0x33, 0x88, 0x42, 0x50,
		0xE3, 0x24, 0x8C, 0x65, 0x7D, 0x41, 0x3B, 0x8C,
		0xE0, 0x1C, 0x8C, 0x9A, 0xAD, 0xED, 0xF5, 0xB9,
	}
	iv := []byte{
		0x9D, 0xEA, 0xDE, 0xC2, 0x62, 0x17, 0x47, 0xA6,
		0x2A, 0x80, 0xA7, 0xC3, 0xFF, 0xA8, 0xE3, 0x47,
	}

	plaintext, err := BeltCFBDecrypt(ciphertext, key, iv)
	if err != nil {
		return fmt.Errorf("ошибка дешифрования: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8text, err := decoder.Bytes(plaintext)
	if err != nil {
		return fmt.Errorf("декодирование Windows-1251: %w", err)
	}

	err = os.WriteFile(outputPath, utf8text, 0644)
	if err != nil {
		return fmt.Errorf("запись в %s: %w", outputPath, err)
	}

	fmt.Println("Расшифрованный текст сохранён в", outputPath)
	return nil
}
