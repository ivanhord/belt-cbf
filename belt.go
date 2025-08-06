package main

/*
#cgo CFLAGS: -I./bee2/include
#cgo LDFLAGS: -L./bee2/build-gcc/src -lbee2_static

#include "bee2/crypto/belt.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func BeltHashGo(data []byte) ([]byte, error) {
	const hashSize = 32 // 256 бит

	out := make([]byte, hashSize)
	ret := C.beltHash(
		(*C.uchar)(unsafe.Pointer(&out[0])),
		unsafe.Pointer(&data[0]),
		C.size_t(len(data)),
	)

	if ret != 0 {
		return nil, errors.New("beltHash error")
	}

	return out, nil
}

func BeltCFBDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	if len(ciphertext) == 0 || len(key) == 0 || len(iv) == 0 {
		return nil, errors.New("некорректные параметры дешифрования")
	}

	plaintext := make([]byte, len(ciphertext))

	ret := C.beltCFBDecr(
		unsafe.Pointer(&plaintext[0]),
		unsafe.Pointer(&ciphertext[0]),
		C.size_t(len(ciphertext)),
		(*C.uchar)(unsafe.Pointer(&key[0])),
		C.size_t(len(key)),
		(*C.uchar)(unsafe.Pointer(&iv[0])),
	)

	if ret != 0 {
		return nil, errors.New("beltCFBDecr error")
	}

	return plaintext, nil
}
