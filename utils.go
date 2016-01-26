package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math"
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

/*
	Round округляет значение val. Возвращает округлённое значение.
	Параметр roundOn задаёт значение разряда, по которому определяется
	вид округления - в большую или в меньшую сторону.
	Параметр places определяет количество знаков после десятичной точки,
	в случае, если он положителен, до целых - если 0. Может быть отрицательным,
	в этом случае, например, при -1 округление выполняется до десятков.
	Примеры:
	round(2.34, .5, 1) возвращает 2.3
	round(2.37, .5, 1) возвращает 2.4
	round(2.37, .5, 0) возвращает 2.0
	round(2.77, .5, 0) возвращает 3.0
*/
func Round(val float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	res := round / pow
	return res
}
