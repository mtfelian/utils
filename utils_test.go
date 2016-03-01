package utils

import (
	"fmt"
	"math"
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

	_, err = UniqID(0)
	if err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

type roundTestCase struct {
	val     float64
	roundOn float64
	places  int
	res     float64
}

func TestRound(t *testing.T) {
	testCases := []roundTestCase{
		{2.34, .5, 1, 2.3},
		{2.37, .5, 1, 2.4},
		{2.37, .5, 0, 2.0},
		{2.77, .5, 0, 3.0},
		{-2.34, .5, 1, -2.3},
		{-2.37, .5, 1, -2.4},
		{-2.37, .5, 0, -2.0},
		{-2.77, .5, 0, -3.0},
		{2.44, .5, 1, 2.4},
		{2.45, .5, 1, 2.5},
		{2.46, .5, 1, 2.5},
		{2.42, .3, 1, 2.4},
		{2.43, .3, 1, 2.5},
		{2.44, .3, 1, 2.5},
		{2.22, .3, 0, 2.0},
		{3.33, .3, 0, 4.0},
		{4.44, .3, 0, 5.0},
		{-2.22, .3, 0, -2.0},
		{-3.33, .3, 0, -4.0},
		{-4.44, .3, 0, -5.0},
		{2.22, .0, 0, 3.0},
		{3.33, .0, 0, 4.0},
		{4.44, .0, 0, 5.0},
		{2.22, .0, 1, 2.3},
		{3.33, .0, 1, 3.4},
		{4.44, .0, 1, 4.5},
		{24.4, .5, -1, 20.0},
		{25.5, .5, -1, 30.0},
		{26.6, .5, -1, 30.0},
	}

	for i, v := range testCases {
		r := Round(v.val, v.roundOn, v.places)
		if math.IsNaN(v.res) && math.IsNaN(r) {
			continue
		}

		if r != v.res {
			t.Fatalf("Функция Round() вернула %v на кейсе %v с индексом %d", r, v, i)
		}
	}
}

func TestFormatPhone(t *testing.T) {
	testCases := []struct {
		input  string
		output string
		err    bool
	}{
		{"+71234567890", "71234567890", false},
		{"71234567890", "71234567890", false},
		{"+7 (861) 12-12-123", "78611212123", false},
		{"8 (861) 12-12-123", "78611212123", false},
		{"8611212123", "8611212123", true},
		{"1234567", "1234567", true},
	}

	for _, testCase := range testCases {
		receivedOutput, err := formatPhone(testCase.input)
		if (err != nil) != testCase.err {
			t.Fatalf("Неверное состояние ошибки на кортеже %v", testCase)
		}
		if receivedOutput != testCase.output {
			t.Fatalf("Полученные данные (%s) не соответствуют ожидаемым (%s)",
				receivedOutput, testCase.output)
		}
	}

}