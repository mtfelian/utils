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
