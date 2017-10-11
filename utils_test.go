package utils

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"fmt"
	"github.com/kr/pretty"
	"time"
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

// TestSliceContains checks checking element in slice
func TestSliceContains(t *testing.T) {
	type testDataElement struct {
		needle   interface{} // what to search
		haystack interface{} // where to search, should be Kind() of reflect.Slice
		result   bool        // expected result
	}

	testData := []testDataElement{
		{interface{}(2), interface{}([]int{1, 2, 3, 4, 5}), true},
		{interface{}(2), interface{}([]int{1, 3, 4, 5}), false},
		{interface{}(2), interface{}([]int{}), false},
		{interface{}(2), interface{}([]int{2}), true},
		{interface{}(2), interface{}([]string{"2"}), false},
		{interface{}(uint(2)), interface{}([]uint{1, 2, 3, 4, 5}), true},
		{interface{}(uint(2)), interface{}([]uint{1, 3, 4, 5}), false},
		{interface{}(uint(2)), interface{}([]uint{}), false},
		{interface{}(uint(2)), interface{}([]uint{2}), true},
		{interface{}(uint(2)), interface{}([]int{2}), false},
		{interface{}(2), interface{}([]uint{2}), false},
		{interface{}("2"), interface{}([]string{"1", "2", "3", "4", "5"}), true},
		{interface{}("2"), interface{}([]string{"1", "3", "4", "5"}), false},
		{interface{}("2"), interface{}([]string{}), false},
		{interface{}("2"), interface{}([]string{"2"}), true},
	}

	for i, value := range testData {
		receivedResult := SliceContains(value.needle, value.haystack)
		if receivedResult != value.result {
			t.Fatalf("Index %d. Expected %v, received %v.", i, value.result, receivedResult)
		}
	}
}

// TestCircularAdd checks circular addition
func TestCircularAdd(t *testing.T) {
	type testDataElement struct {
		a      int
		max    int
		result int
	}

	testData := []testDataElement{
		{0, 0, 0},
		{0, 1, 1},
		{1, 1, 0},
		{0, 2, 1},
		{1, 2, 2},
		{2, 2, 0},
	}

	for i, value := range testData {
		receivedResult := CircularAdd(value.a, value.max)
		if receivedResult != value.result {
			t.Fatalf("Index %d. CircularAdd(%d, %d). Expected %d, received %d.",
				i, value.a, value.max, value.result, receivedResult)
		}
	}
}

// TestIndicesSlice checks IndicesSlice sorting
func TestIndicesSlice(t *testing.T) {
	type testDataElement struct {
		sourceSlice     []uint
		expectedSlice   UintSlice
		expectedIndices []int // indices from source slice after sorting
	}

	testData := []testDataElement{
		{[]uint{6, 2, 1, 4, 3}, UintSlice{1, 2, 3, 4, 6}, []int{2, 1, 4, 3, 0}},
		{[]uint{2, 1, 3}, UintSlice{1, 2, 3}, []int{1, 0, 2}},
		{[]uint{2}, UintSlice{2}, []int{0}},
		{[]uint{}, UintSlice{}, []int{}},
		{[]uint{59, 23, 1, 3, 23}, UintSlice{1, 3, 23, 23, 59}, []int{2, 3, 1, 4, 0}},
	}

	for i, value := range testData {
		slice := NewIndicesUintSlice(value.sourceSlice...)
		sort.Sort(slice)
		if !reflect.DeepEqual(slice.Indices, value.expectedIndices) {
			t.Fatalf("Index %d, wrong indices: %s", i, pretty.Diff(slice.Indices, value.expectedIndices))
		}

		underlyingSlice, ok := slice.Interface.(UintSlice)
		if !ok {
			t.Fatalf("Can't get underlying slice of slice.Interface")
		}

		if !reflect.DeepEqual(underlyingSlice, value.expectedSlice) {
			t.Fatalf("Index %d, wrong received slice: %s", i, pretty.Diff(underlyingSlice, value.expectedSlice))
		}
	}
}

// TestStringToUintSlice checks converting sepatated string to slice of uint values
func TestStringToUintSlice(t *testing.T) {
	type testDataElement struct {
		sourceString  string
		separator     string
		min           uint
		expectedSlice []uint
		expectedError bool
	}

	testData := []testDataElement{
		{"4,5,6", ",", 2, []uint{4, 5, 6}, false},   // 0
		{"4|5|6", "|", 5, []uint{5, 6}, false},      // 1
		{",4,5,6,", ",", 2, []uint{4, 5, 6}, false}, // 2
		{"10", ",", 10, []uint{10}, false},          // 3
		{"10,0", ",", 10, []uint{10}, false},        // 4
		{"10,q", ",", 10, []uint{}, true},           // 5
		{",", ",", 10, []uint{}, false},             // 6
		{"", ",", 10, []uint{}, false},              // 7
	}

	for i, value := range testData {
		receivedSlice, err := StringToUintSlice(value.sourceString, value.separator, value.min)
		if (err != nil) != value.expectedError {
			t.Fatalf("Index %d, expected is error: %v, received: %v", i, value.expectedError, err != nil)
		}

		if !reflect.DeepEqual(value.expectedSlice, receivedSlice) {
			t.Fatalf("Index %d, wrong received slice, diff: %s",
				i, pretty.Diff(receivedSlice, value.expectedSlice))
		}
	}
}

// TestStringToStringSlice checks converting sepatated string to slice of string values
func TestStringToStringSlice(t *testing.T) {
	type testDataElement struct {
		sourceString  string
		separator     string
		expectedSlice []string
	}

	testData := []testDataElement{
		{"q4,5,6", ",", []string{"q4", "5", "6"}}, // 0
		{",4,5,6,", ",", []string{"4", "5", "6"}}, // 1
		{"10", ",", []string{"10"}},               // 2
		{"10,0", ",", []string{"10", "0"}},        // 3
		{"10,q", ",", []string{"10", "q"}},        // 4
		{",", ",", []string{}},                    // 5
		{"", ",", []string{}},                     // 6
	}

	for i, value := range testData {
		receivedSlice := StringToStringSlice(value.sourceString, value.separator)
		if !reflect.DeepEqual(value.expectedSlice, receivedSlice) {
			t.Fatalf("Index %d, wrong received slice, diff: %s",
				i, pretty.Diff(receivedSlice, value.expectedSlice))
		}
	}
}

// TestToLowerFirstRune checks converting first rune of string to lowercase
func TestToLowerFirstRune(t *testing.T) {
	type testDataElement struct {
		sourceString   string
		expectedString string
	}

	testData := []testDataElement{
		{"", ""},             // 0
		{"Q", "q"},           // 1
		{"q", "q"},           // 2
		{":", ":"},           // 3
		{" ", " "},           // 4
		{"QWERTY", "qWERTY"}, // 5
		{"qwerty", "qwerty"}, // 6
		{"Q Q", "q Q"},       // 7
		{"Q_Q", "q_Q"},       // 8
	}

	for i, value := range testData {
		receivedString := ToLowerFirstRune(value.sourceString)
		if receivedString != value.expectedString {
			t.Fatalf("Index %d, expected: %s, received: %s", i, value.expectedString, receivedString)
		}
	}
}

// TestTry checks trying function
func TestTry(t *testing.T) {
	i := 0
	actionFunc := func() error {
		if i < 5 {
			fmt.Printf("success: executed actionFunc() with i = %d\n", i)
			return nil
		}
		fmt.Printf("fail: executed actionFunc() with i = %d\n", i)
		return fmt.Errorf("Error. i too high: %d", i)
	}
	conditionFunc := func(e error) bool { return e != nil }
	for i < 6 {
		i++
		shouldAttempt := 2
		attempts, err := Try(actionFunc, shouldAttempt, time.Millisecond, conditionFunc)
		if i < 5 {
			if attempts != 1 {
				t.Fatalf("Expected attempts: %d, received: %d", 1, attempts)
			}
		} else {
			if err == nil {
				t.Fatalf("Expected error, received nil")
			}
			if attempts != shouldAttempt {
				t.Fatalf("Expected attempts: %d, received: %d", shouldAttempt, attempts)
			}
		}
	}
}
