package main

import (
	"belt-cbf/internal/bee2"
	"encoding/hex"
	"errors"
	"strings"
)

func DecryptHex(hexCipher string) ([]byte, error) {
	hexCipher = strings.TrimSpace(hexCipher)

	// Удаляем префиксы
	for _, p := range []string{"0x", "0X", "\\x"} {
		if strings.HasPrefix(hexCipher, p) {
			hexCipher = strings.TrimPrefix(hexCipher, p)
			break
		}
	}

	ciphertext, err := hex.DecodeString(hexCipher)
	if err != nil {
		return nil, errors.New("ошибка декодирования HEX: " + err.Error())
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

	plaintext, err := bee2.BeltCFBDecrypt(ciphertext, key, iv)
	if err != nil {
		return nil, errors.New("ошибка дешифрования: " + err.Error())
	}

	return plaintext, nil
}
