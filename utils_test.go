package utils

import (
	"fmt"
	"testing"
)

func TestStringToUint(t *testing.T) {
	testInt := uint(8)
	testString := fmt.Sprintf("%d", testInt)
	res, err := StringToUint(testString)
	if err != nil {
		t.Fatal(err)
	}
	if res != testInt {
		t.Fatalf("Ожидалось получить %d, получено %d", testInt, res)
	}
}

func TestBadStringToUint(t *testing.T) {
	testString := "q"
	_, err := StringToUint(testString)
	if err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestUniqID(t *testing.T) {
	str, err := UniqID(13)
	if err != nil {
		t.Fatal(err)
	}
	if len(str) != 13 {
		t.Fatalf("Длина строки не 13, а %d", len(str))
	}

	str, err = UniqID(0)
	if err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}
