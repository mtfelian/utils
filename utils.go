package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math"
	"os"
	"path/filepath"
	"regexp"
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

// FormatPhone форматирует строку с номером телефона в формат "71234567890"
// Возвращает:
// Успех: Форматированный номер телефона, nil
// Ошибка: Исходный номер телефона, ошибка
func FormatPhone(phone string) (string, error) {
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

var renderFloatPrecisionMultipliers = [10]float64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
}

var renderFloatPrecisionRounders = [10]float64{
	0.5,
	0.05,
	0.005,
	0.0005,
	0.00005,
	0.000005,
	0.0000005,
	0.00000005,
	0.000000005,
	0.0000000005,
}

// todo документация
func RenderFloat(format string, n float64) (string, error) {
	if math.IsNaN(n) {
		return "NaN", nil
	}
	if n > math.MaxFloat64 {
		return "Infinity", nil
	}
	if n < -math.MaxFloat64 {
		return "-Infinity", nil
	}

	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		precision = 9
		thousandStr = ""

		formatDirectiveChars := []rune(format)
		formatDirectiveIndices := make([]int, 0)
		for i, char := range formatDirectiveChars {
			if char != '#' && char != '0' {
				formatDirectiveIndices = append(formatDirectiveIndices, i)
			}
		}

		if len(formatDirectiveIndices) > 0 {
			if formatDirectiveIndices[0] == 0 {
				if formatDirectiveChars[formatDirectiveIndices[0]] != '+' {
					return "", errors.New("RenderFloat(): ошибка, должен быть положительный знак")
				}
				positiveStr = "+"
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			if len(formatDirectiveIndices) == 2 {
				if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
					return "", errors.New("RenderFloat(): ошибка, за разделителем разрядов тысяч должны следовать три спецификатора цифры")
				}
				thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			if len(formatDirectiveIndices) == 1 {
				decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
			}
		}
	}

	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	intStr := strconv.Itoa(int(intf))

	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	if precision == 0 {
		return signStr + intStr, nil
	}

	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr, nil
}

// todo документация
func RenderInteger(format string, n int64) (string, error) {
	return RenderFloat(format, float64(n))
}

// GetSelfPath возвращает путь к исполняемому файлу. Или ошибку
func GetSelfPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// FileExists возвращает true если файл с путём существует. Иначе false
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir возвращает true если указанный путь является директорией, иначе false
func IsDir(path string) (bool, error) {
	if !FileExists(path) {
		return false, errors.New("Файл не существует")
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}
