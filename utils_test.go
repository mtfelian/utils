package utils

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"testing"
)

type testStringToUintResult struct {
	Value uint
	OK    bool
}

func TestStringToUint(t *testing.T) {
	testCases := map[string]testStringToUintResult{
		"0": testStringToUintResult{Value: 0, OK: true},
		"7": testStringToUintResult{Value: 7, OK: true},
		"q": testStringToUintResult{Value: 0, OK: false},
		"":  testStringToUintResult{Value: 0, OK: false},
	}

	for param, testCase := range testCases {
		receivedValue, err := StringToUint(param)
		if (err == nil) != testCase.OK || receivedValue != testCase.Value {
			t.Fatalf("Тест не пройден на наборе %s:%#v", param, testCase)
		}
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
		receivedOutput, err := FormatPhone(testCase.input)
		if (err != nil) != testCase.err {
			t.Fatalf("Неверное состояние ошибки на кортеже %v", testCase)
		}
		if receivedOutput != testCase.output {
			t.Fatalf("Полученные данные (%s) не соответствуют ожидаемым (%s)",
				receivedOutput, testCase.output)
		}
	}

}

func TestFileExists(t *testing.T) {
	binPath, err := GetSelfPath()
	if err != nil {
		t.Fatal(err)
	}

	fileName := filepath.Join(binPath, "test.txt")
	defer os.Remove(fileName)

	err = ioutil.WriteFile(fileName, []byte("test"), 0660)
	if err != nil {
		t.Fatalf("Ошибка записи файла: %v", err)
	}

	if !FileExists(fileName) {
		t.Fatal("Но файл существует!")
	}

	err = os.Remove(fileName)
	if err != nil {
		t.Fatalf("Ошибка удаления файла: %v", err)
	}

	if FileExists(fileName) {
		t.Fatal("Но файл не существует!")
	}
}

func TestIsDirFalse(t *testing.T) {
	binPath, err := GetSelfPath()
	if err != nil {
		t.Fatal(err)
	}

	fileName := filepath.Join(binPath, "test.txt")
	defer os.Remove(fileName)

	err = ioutil.WriteFile(fileName, []byte("test"), 0660)
	if err != nil {
		t.Fatalf("Ошибка записи файла: %v", err)
	}

	isDir, err := IsDir(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if isDir {
		t.Fatal("Ожидалось получить isDir false, получили true")
	}
}

func TestIsDirTrue(t *testing.T) {
	binPath, err := GetSelfPath()
	if err != nil {
		t.Fatal(err)
	}

	dirName := filepath.Join(binPath, "testDir")
	defer os.RemoveAll(dirName)
	if err := os.Mkdir(dirName, 0777); err != nil {
		t.Fatal(err)
	}

	isDir, err := IsDir(dirName)
	if err != nil {
		t.Fatal(err)
	}
	if !isDir {
		t.Fatal("Ожидалось получить isDir true, получили false")
	}
}

func TestEncodeToWindows1251(t *testing.T) {
	str := "Я проверяю tochno"
	strToWindows1251, err := EncodeToWindows1251([]byte(str))
	if err != nil {
		t.Fatal(err)
	}
	if len(strToWindows1251) != 17 {
		t.Fatal("Ошибка преобразования в windows-1251, русские буквы должны занимать 1 байт")
	}
}

func TestTrimSnils(t *testing.T) {
	testCases := map[string]string{
		"12345678910":      "12345678910",
		"123-456- 789-1 0": "12345678910",
		"123 ":             "123",
		"":                 "",
	}

	for key, expectedValue := range testCases {
		receivedValue := trimSnils(key)
		if expectedValue != receivedValue {
			t.Fatalf("Кейс '%s'. Ожидалось '%s', получено '%s'", key, expectedValue, receivedValue)
		}
	}
}

func TestCheckSnils(t *testing.T) {
	type validationResult struct {
		ok  bool
		err error
	}
	// err не nil хз как тут протестить
	testCases := map[string]validationResult{
		"13972606386":    validationResult{ok: true, err: nil},
		"16776206804":    validationResult{ok: true, err: nil},
		"167 762 068-04": validationResult{ok: true, err: nil},
		"1677620680":     validationResult{ok: false, err: nil},
		"00000000000":    validationResult{ok: false, err: nil},
		"00000000100":    validationResult{ok: true, err: nil},
		"00100199799":    validationResult{ok: true, err: nil},
		"00100199899":    validationResult{ok: false, err: nil},
		"10050010056":    validationResult{ok: false, err: nil},
		"167762o6804":    validationResult{ok: false, err: nil},
	}

	for key, expectedValue := range testCases {
		receivedValue, receivedErr := CheckSnils(key)
		if (expectedValue.err == nil) != (receivedErr == nil) {
			t.Fatalf("Кейс '%s'. Err ожидалось '%v', получено '%v'",
				key, expectedValue.err, receivedErr)
		}
		if expectedValue.ok != receivedValue {
			t.Fatalf("Кейс '%s'. Проверка ожидалась '%v', получено '%v'",
				key, expectedValue.ok, receivedValue)
		}
	}
}

func TestCountElementsOnPage(t *testing.T) {
	// {elementsTotal, page, pageSize}: result
	testData := map[[3]uint]uint{
		[3]uint{0, 1, 1}:   0,
		[3]uint{1, 1, 1}:   1,
		[3]uint{0, 0, 1}:   0,
		[3]uint{1, 1, 0}:   1,
		[3]uint{1, 2, 1}:   0,
		[3]uint{10, 1, 1}:  1,
		[3]uint{10, 1, 0}:  10,
		[3]uint{10, 0, 1}:  1,
		[3]uint{10, 5, 1}:  1,
		[3]uint{10, 10, 1}: 1,
		[3]uint{10, 11, 1}: 0,
		[3]uint{10, 5, 0}:  0,
		[3]uint{10, 10, 0}: 0,
		[3]uint{10, 11, 0}: 0,
		[3]uint{10, 1, 4}:  4,
		[3]uint{10, 0, 4}:  4,
		[3]uint{10, 2, 4}:  4,
		[3]uint{10, 3, 4}:  2,
		[3]uint{10, 1, 20}: 10,
		[3]uint{10, 0, 20}: 10,
		[3]uint{10, 2, 20}: 0,
	}

	for key, value := range testData {
		receivedResult := CountElementsOnPage(key[0], key[1], key[2])
		if receivedResult != value {
			t.Fatalf("Кортеж %v. Ожидался результат %d, получен результат %d.", key, value, receivedResult)
		}
	}
}

func TestCountPages(t *testing.T) {
	// {elementsTotal, pageSize}: result
	testData := map[[2]uint]uint{
		[2]uint{0, 0}:   0,
		[2]uint{0, 1}:   0,
		[2]uint{1, 0}:   1,
		[2]uint{1, 1}:   1,
		[2]uint{10, 1}:  10,
		[2]uint{10, 5}:  2,
		[2]uint{10, 3}:  4,
		[2]uint{10, 8}:  2,
		[2]uint{10, 0}:  1,
		[2]uint{10, 10}: 1,
		[2]uint{10, 15}: 1,
	}

	for key, value := range testData {
		receivedResult := CountPages(key[0], key[1])
		if receivedResult != value {
			t.Fatalf("Кортеж %v. Ожидался результат %d, получен результат %d.", key, value, receivedResult)
		}
	}
}

func TestStringsToInterfaces(t *testing.T) {
	strings := []string{`123`, `234`, `345`, `123987`, ``}
	interfaces := StringsToInterfaces(strings)
	for i, I := range interfaces {
		backString, valid := I.(string)
		if !valid {
			t.Fatalf("Error converting string at index %d: invalid", i)
		}
		if backString != strings[i] {
			t.Fatalf("Expected: %s, received: %s", strings[i], backString)
		}
	}
}
