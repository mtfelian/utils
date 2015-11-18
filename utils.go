package utils

import (
	"strconv"
)

/*
	Трактует строку s как значение типа uint
*/
func StringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
