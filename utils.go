package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math"
	"strconv"
	"regexp"
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

// Round округляет значение val. Возвращает округлённое значение.
// Параметр roundOn задаёт значение разряда, по которому определяется
// вид округления - в большую или в меньшую сторону.
// Параметр places определяет количество знаков после десятичной точки,
// в случае, если он положителен, до целых - если 0. Может быть отрицательным,
// в этом случае, например, при -1 округление выполняется до десятков.
// Примеры:
// Round(2.34, .5, 1) возвращает 2.3
// Round(2.37, .5, 1) возвращает 2.4
// Round(2.37, .5, 0) возвращает 2.0
// Round(2.77, .5, 0) возвращает 3.0
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

// formatPhone форматирует строку с номером телефона в формат "71234567890"
// Возвращает:
// Успех: Форматированный номер телефона, nil
// Ошибка: Исходный номер телефона, ошибка
func formatPhone(phone string) (string, error) {
	// форматируем строку с телефоном
	res := phone
	reg, err := regexp.Compile(`[\(\).,;#*А-яA-z\s+-]*`)
	if err != nil {
		return phone, err
	}
	res = reg.ReplaceAllString(phone, "")
	// длина строки с телефоном в норме должна быть 12 символов если с "+" или 11 символов без оного
	if len(res) > 11 {
		return phone, errors.New("Слишком длинный номер телефона")
	} else if len(res) < 11 {
		return phone, errors.New("Слишком короткий номер телефона")
	}
	if res[:1] == "8" {
		res = "7" + res[1:len(res)]
	}
	return res, nil
}