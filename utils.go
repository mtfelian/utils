package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
)

// StringToUint трактует строку s как значение типа uint
func StringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// UniqID формирует уникальную строку длины n
func UniqID(n int) (string, error) {
	if n < 1 {
		return "", errors.New("n должно быть больше 0")
	}
	b := make([]byte, n)
	// здесь err == nil только если мы читаем len(b) байт
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	result := base64.URLEncoding.EncodeToString(b)
	if len(result) > n {
		return result[:n], nil
	}
	return result, nil
}
