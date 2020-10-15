package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with Ginkgo", func() {

	It("checks StringToUint", func() {
		type testStringToUintResult struct {
			Value uint
			OK    bool
		}

		testCases := map[string]testStringToUintResult{
			"0": {Value: 0, OK: true},
			"7": {Value: 7, OK: true},
			"q": {Value: 0, OK: false},
			"":  {Value: 0, OK: false},
		}

		for param, testCase := range testCases {
			By(fmt.Sprintf("testing on %s:%#v", param, testCase))
			receivedValue, err := StringToUint(param)
			Expect((err == nil) != testCase.OK || receivedValue != testCase.Value).To(BeFalse())
		}
	})

	It("checks UniqID", func() {
		str, err := UniqID(13)
		Expect(err).NotTo(HaveOccurred())
		Expect(str).To(HaveLen(13))

		_, err = UniqID(0)
		Expect(err).To(HaveOccurred())
	})

	It("checks Round", func() {
		testCases := []struct {
			val     float64
			roundOn float64
			places  int
			res     float64
		}{
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
			By(fmt.Sprintf("testing case %d, %v", i, v))
			r := Round(v.val, v.roundOn, v.places)
			if math.IsNaN(v.res) && math.IsNaN(r) {
				continue
			}
			Expect(r).To(Equal(v.res))
		}
	})

	It("checks FormatPhone", func() {
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
			By(fmt.Sprintf("testing case %v", testCase))
			receivedOutput, err := FormatPhone(testCase.input)
			Expect(err != nil).To(Equal(testCase.err))
			Expect(receivedOutput).To(Equal(testCase.output))
		}
	})

	It("checks FileExists", func() {
		binPath := MustSelfPath()
		fileName := filepath.Join(binPath, "test.txt")
		defer os.Remove(fileName)

		Expect(ioutil.WriteFile(fileName, []byte("test"), 0660)).To(Succeed())
		Expect(FileExists(fileName)).To(BeTrue())
		Expect(os.Remove(fileName)).To(Succeed())
		Expect(FileExists(fileName)).To(BeFalse())
	})

	It("checks FileSize", func() {
		binPath := MustSelfPath()
		fileName := filepath.Join(binPath, "test.txt")
		defer os.Remove(fileName)

		content := "test\ntest"
		Expect(ioutil.WriteFile(fileName, []byte(content), 0660)).To(Succeed())
		Expect(FileSize(fileName)).To(BeNumerically("==", len(content)))
		Expect(os.Remove(fileName)).To(Succeed())
		Expect(FileSize(fileName)).To(BeNumerically("==", 0))
	})

	Context("IsDir()", func() {
		It("should be false", func() {
			binPath := MustSelfPath()
			fileName := filepath.Join(binPath, "test.txt")
			defer os.Remove(fileName)

			Expect(ioutil.WriteFile(fileName, []byte("test"), 0660)).To(Succeed())
			isDir, err := IsDir(fileName)
			Expect(err).NotTo(HaveOccurred())
			Expect(isDir).To(BeFalse())
		})

		It("should be true", func() {
			binPath := MustSelfPath()

			dirName := filepath.Join(binPath, "testDir")
			defer os.RemoveAll(dirName)
			Expect(os.Mkdir(dirName, 0777)).To(Succeed())

			isDir, err := IsDir(dirName)
			Expect(err).NotTo(HaveOccurred())
			Expect(isDir).To(BeTrue())
		})
	})

	It("checks EncodeToWindows1251", func() {
		str := "Я проверяю tochno"
		strToWindows1251, err := EncodeToWindows1251([]byte(str))
		Expect(err).NotTo(HaveOccurred())
		Expect(strToWindows1251).To(HaveLen(17))
	})

	Context("SNILS", func() {
		It("checks TrimSnils", func() {
			testCases := map[string]string{
				"12345678910":      "12345678910",
				"123-456- 789-1 0": "12345678910",
				"123 ":             "123",
				"":                 "",
			}

			for key, expectedValue := range testCases {
				By(fmt.Sprintf("testing case %s: %s", key, expectedValue))
				Expect(trimSnils(key)).To(Equal(expectedValue))
			}
		})

		It("checks CheckSnils", func() {
			type validationResult struct {
				ok  bool
				err error
			}

			testCases := map[string]validationResult{
				"13972606386":    {ok: true, err: nil},
				"16776206804":    {ok: true, err: nil},
				"167 762 068-04": {ok: true, err: nil},
				"1677620680":     {ok: false, err: nil},
				"00000000000":    {ok: false, err: nil},
				"00000000100":    {ok: true, err: nil},
				"00100199799":    {ok: true, err: nil},
				"00100199899":    {ok: false, err: nil},
				"10050010056":    {ok: false, err: nil},
				"167762o6804":    {ok: false, err: nil},
			}

			for key, expectedValue := range testCases {
				By(fmt.Sprintf("testing case %s: %v", key, expectedValue))
				receivedValue, receivedErr := CheckSnils(key)
				Expect(expectedValue.err == nil).To(Equal(receivedErr == nil))
				Expect(expectedValue.ok).To(Equal(receivedValue))
			}
		})
	})

	It("checks CountElementsOnPage", func() {
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
			By(fmt.Sprintf("testing case %v: %d", key, value))
			Expect(CountElementsOnPage(key[0], key[1], key[2])).To(Equal(value))
		}
	})

	It("checks CountPages", func() {
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
			By(fmt.Sprintf("testing case %v: %d", key, value))
			Expect(CountPages(key[0], key[1])).To(Equal(value))
		}
	})

	It("checks SliceContains", func() {
		testData := []struct {
			needle   interface{}
			haystack interface{}
			result   bool
		}{
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
			{interface{}("2"), interface{}([]string(nil)), false},
		}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d", i))
			Expect(SliceContains(value.needle, value.haystack)).To(Equal(value.result))
		}
	})

	It("checks CircularAdd", func() {
		testData := []struct {
			a      int
			max    int
			result int
		}{{0, 0, 0}, {0, 1, 1}, {1, 1, 0}, {0, 2, 1}, {1, 2, 2}, {2, 2, 0}}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d: %v", i, value))
			Expect(CircularAdd(value.a, value.max)).To(Equal(value.result))
		}
	})

	It("checks NewIndicesUintSlice", func() {
		testData := []struct {
			sourceSlice     []uint
			expectedSlice   UintSlice
			expectedIndices []int
		}{
			{[]uint{6, 2, 1, 4, 3}, UintSlice{1, 2, 3, 4, 6}, []int{2, 1, 4, 3, 0}},
			{[]uint{2, 1, 3}, UintSlice{1, 2, 3}, []int{1, 0, 2}},
			{[]uint{2}, UintSlice{2}, []int{0}},
			{[]uint{}, UintSlice{}, []int{}},
			{[]uint{59, 23, 1, 3, 23}, UintSlice{1, 3, 23, 23, 59}, []int{2, 3, 1, 4, 0}},
		}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d: %v", i, value))
			slice := NewIndicesUintSlice(value.sourceSlice...)
			sort.Sort(slice)
			Expect(slice.Indices).To(Equal(value.expectedIndices))
			underlyingSlice, ok := slice.Interface.(UintSlice)
			Expect(ok).To(BeTrue())
			Expect(underlyingSlice).To(Equal(value.expectedSlice))
		}
	})

	It("checks StringToUintSlice", func() {
		testData := []struct {
			sourceString  string
			separator     string
			min           uint
			expectedSlice []uint
			expectedError bool
		}{
			{"4,5,6", ",", 2, []uint{4, 5, 6}, false},
			{"4|5|6", "|", 5, []uint{5, 6}, false},
			{",4,5,6,", ",", 2, []uint{4, 5, 6}, false},
			{"10", ",", 10, []uint{10}, false},
			{"10,0", ",", 10, []uint{10}, false},
			{"10,q", ",", 10, []uint{}, true},
			{",", ",", 10, []uint{}, false},
			{"", ",", 10, []uint{}, false},
		}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d: %v", i, value))
			receivedSlice, err := StringToUintSlice(value.sourceString, value.separator, value.min)
			Expect(err != nil).To(Equal(value.expectedError))
			Expect(receivedSlice).To(Equal(value.expectedSlice))
		}
	})

	It("checks StringToStringSlice", func() {
		testData := []struct {
			sourceString  string
			separator     string
			expectedSlice []string
		}{
			{"q4,5,6", ",", []string{"q4", "5", "6"}},
			{",4,5,6,", ",", []string{"4", "5", "6"}},
			{"10", ",", []string{"10"}},
			{"10,0", ",", []string{"10", "0"}},
			{"10,q", ",", []string{"10", "q"}},
			{",", ",", []string{}},
			{"", ",", []string{}},
		}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d: %v", i, value))
			Expect(StringToStringSlice(value.sourceString, value.separator)).To(Equal(value.expectedSlice))
		}
	})

	It("checks ToLowerFirstRune", func() {
		testData := []struct {
			sourceString   string
			expectedString string
		}{
			{"", ""}, {"Q", "q"}, {"q", "q"}, {":", ":"}, {" ", " "},
			{"QWERTY", "qWERTY"}, {"qwerty", "qwerty"}, {"Q Q", "q Q"}, {"Q_Q", "q_Q"},
		}

		for i, value := range testData {
			By(fmt.Sprintf("testing case %d: %v", i, value))
			Expect(ToLowerFirstRune(value.sourceString)).To(Equal(value.expectedString))
		}
	})

	It("checks Try", func() {
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
				Expect(attempts).To(Equal(1))
			} else {
				Expect(err).To(HaveOccurred())
				Expect(attempts).To(Equal(shouldAttempt))
			}
		}
	})

})

var _ = Describe("Between func", func() {
	It("works", func() {
		testCases := []struct {
			source string
			l, r   string
			result string
		}{
			{
				source: "12345687890",
				l:      "2", r: "45",
				result: "3",
			},
			{
				source: "12345687890",
				l:      "45", r: "45",
				result: "",
			},
			{
				source: "12345687890",
				l:      "45", r: "2",
				result: "",
			},
			{
				source: "11111111",
				l:      "1", r: "11",
				result: "",
			},
			{
				source: "111222333111",
				l:      "111", r: "333",
				result: "222",
			},
			{
				source: "111222333111",
				l:      "222", r: "111",
				result: "333",
			},
			{
				source: "",
				l:      "111", r: "333",
				result: "",
			},
			{
				source: "",
				l:      "", r: "",
				result: "",
			},
		}

		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, source: %s", i, tc.source))
			Expect(SubstringBetween(tc.source, tc.l, tc.r)).To(Equal(tc.result))
		}
	})
})

var _ = Describe("BackupFileName func", func() {
	It("works", func() {
		testCases := []struct {
			createFiles       []string
			inputFileName     string
			inputExtension    string
			resultingFileName string
		}{
			{
				createFiles:       []string{"a"},
				inputFileName:     "a",
				inputExtension:    "bak",
				resultingFileName: "a.bak1",
			},
			{
				createFiles:       []string{"a", "a.bak1"},
				inputFileName:     "a",
				inputExtension:    "bak",
				resultingFileName: "a.bak2",
			},
			{
				createFiles:       []string{"a", "a.bak"},
				inputFileName:     "a.bak",
				inputExtension:    "bak",
				resultingFileName: "a.bak1",
			},
			{
				createFiles:       []string{"a", "a.bak", "a.bak1"},
				inputFileName:     "a.bak1",
				inputExtension:    "bak",
				resultingFileName: "a.bak2",
			},
			{
				createFiles:       []string{"a", "a.bak"},
				inputFileName:     "a.bak",
				inputExtension:    "",
				resultingFileName: "a.bak1",
			},
			{
				createFiles:       []string{"a", "a.bak", "a.bak1"},
				inputFileName:     "a.bak1",
				inputExtension:    "ext",
				resultingFileName: "a.bak1.ext1",
			},
			{
				createFiles:       []string{"a", "a.bak"},
				inputFileName:     "a.bak",
				inputExtension:    "bak",
				resultingFileName: "a.bak1",
			},
		}

		selfPath := MustSelfPath()
		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, inputFileName: %s", i, tc.inputFileName))
			// creating files
			for _, fn := range tc.createFiles {
				Expect(ioutil.WriteFile(filepath.Join(selfPath, fn), []byte{}, 0660)).To(Succeed())
			}
			Expect(BackupFileName(filepath.Join(selfPath, tc.inputFileName), tc.inputExtension)).
				To(Equal(filepath.Join(selfPath, tc.resultingFileName)))
			// removing files
			for _, fn := range tc.createFiles {
				Expect(os.Remove(filepath.Join(selfPath, fn))).To(Succeed())
			}
		}
	})
})

var _ = Describe("RemoveDuplicates func", func() {
	It("works for []uint", func() {
		testCases := []struct {
			slice          []uint
			expectedResult []uint
		}{
			{slice: []uint{}, expectedResult: []uint{}},
			{slice: []uint{1}, expectedResult: []uint{1}},
			{slice: []uint{2, 1}, expectedResult: []uint{2, 1}},
			{slice: []uint{2, 2}, expectedResult: []uint{2}},
			{slice: []uint{3, 2, 2, 3, 2}, expectedResult: []uint{3, 2}},
			{slice: []uint{2, 3, 3, 2, 3}, expectedResult: []uint{2, 3}},
		}
		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, slice: %v", i, tc.slice))
			result, hadNaN, err := RemoveDuplicates(tc.slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(hadNaN).To(BeFalse())
			Expect(result.([]uint)).To(Equal(tc.expectedResult))
		}
	})

	It("works for []int", func() {
		testCases := []struct {
			slice          []int
			expectedResult []int
		}{
			{slice: []int{}, expectedResult: []int{}},
			{slice: []int{1}, expectedResult: []int{1}},
			{slice: []int{2, 1}, expectedResult: []int{2, 1}},
			{slice: []int{2, 2}, expectedResult: []int{2}},
			{slice: []int{3, 2, 2, 3, 2}, expectedResult: []int{3, 2}},
			{slice: []int{2, 3, 3, 2, 3}, expectedResult: []int{2, 3}},
		}
		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, slice: %v", i, tc.slice))
			result, hadNaN, err := RemoveDuplicates(tc.slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(hadNaN).To(BeFalse())
			Expect(result.([]int)).To(Equal(tc.expectedResult))
		}
	})

	It("works for []string", func() {
		testCases := []struct {
			slice          []string
			expectedResult []string
		}{
			{slice: []string{}, expectedResult: []string{}},
			{slice: []string{"1"}, expectedResult: []string{"1"}},
			{slice: []string{"2", "1"}, expectedResult: []string{"2", "1"}},
			{slice: []string{"2", "2"}, expectedResult: []string{"2"}},
			{slice: []string{" 2", "2", "2 "}, expectedResult: []string{" 2", "2", "2 "}},
			{slice: []string{"3", "2", "2", "3", "2"}, expectedResult: []string{"3", "2"}},
			{slice: []string{"2", "3", "3", "2", "3"}, expectedResult: []string{"2", "3"}},
		}
		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, slice: %v", i, tc.slice))
			result, hadNaN, err := RemoveDuplicates(tc.slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(hadNaN).To(BeFalse())
			Expect(result.([]string)).To(Equal(tc.expectedResult))
		}
	})

	It("works for struct", func() {
		type testStruct struct {
			I int
			S string
		}
		testCases := []struct {
			slice          []testStruct
			expectedResult []testStruct
		}{
			{slice: []testStruct{},
				expectedResult: []testStruct{}}, // 0
			{slice: []testStruct{{I: 1, S: "1"}},
				expectedResult: []testStruct{{I: 1, S: "1"}}}, // 1
			{slice: []testStruct{{I: 1, S: "1"}, {I: 2, S: "2"}},
				expectedResult: []testStruct{{I: 1, S: "1"}, {I: 2, S: "2"}}},
			{slice: []testStruct{{I: 2, S: "2"}, {I: 2, S: "2"}}, // 2
				expectedResult: []testStruct{{I: 2, S: "2"}}},
			{slice: []testStruct{{I: 2, S: "2"}, {I: 2, S: "2 "}, {I: 2, S: " 2"}}, // 3
				expectedResult: []testStruct{{I: 2, S: "2"}, {I: 2, S: "2 "}, {I: 2, S: " 2"}}},
			{slice: []testStruct{{I: 3, S: "3"}, {I: 2, S: "2"}, {I: 2, S: "2"}, {I: 3, S: "3"}, {I: 2, S: "2"}}, // 4
				expectedResult: []testStruct{{I: 3, S: "3"}, {I: 2, S: "2"}}}, // 5
			{slice: []testStruct{{I: 2, S: "2"}, {I: 3, S: "3"}, {I: 3, S: "3"}, {I: 2, S: "2"}, {I: 3, S: "3"}},
				expectedResult: []testStruct{{I: 2, S: "2"}, {I: 3, S: "3"}}}, // 6
		}
		for i, tc := range testCases {
			By(fmt.Sprintf("testing %d case, slice: %v", i, tc.slice))
			result, hadNaN, err := RemoveDuplicates(tc.slice)
			Expect(err).NotTo(HaveOccurred())
			Expect(hadNaN).To(BeFalse())
			typedResult := result.([]testStruct)
			Expect(reflect.DeepEqual(typedResult, tc.expectedResult)).To(BeTrue())
		}
	})
})

var _ = Describe("MarshalUnmarshalJSON func", func() {
	It("checks it", func() {
		type dataNested struct {
			NestedValue1 string `json:"nv1"`
			NestedValue2 int    `json:"nv2"`
		}
		type data struct {
			Key1 string     `json:"key1"`
			Key2 dataNested `json:"key2"`
		}
		for i, tc := range []struct {
			jsonData       interface{}
			expectedResult data
			err            bool
		}{
			{
				jsonData: map[string]interface{}{
					"key1": "value",
					"key2": map[string]interface{}{"nv1": "v1", "nv2": 7},
				},
				expectedResult: data{Key1: "value", Key2: dataNested{NestedValue1: "v1", NestedValue2: 7}},
			}, // 0
			{jsonData: nil, expectedResult: data{}},          // 1
			{jsonData: 7, expectedResult: data{}, err: true}, // 2
			{
				jsonData:       map[string]interface{}{"key1": "value", "key2": 7},
				expectedResult: data{Key1: "value"},
				err:            true,
			}, // 3
			{
				jsonData:       map[string]interface{}{"key1": "value", "key2": nil},
				expectedResult: data{Key1: "value"},
			}, // 4
		} {
			By(fmt.Sprintf("testing %d case, jsonData: %v", i, tc.jsonData))
			var out data

			Expect(MarshalUnmarshalJSON(tc.jsonData, &out) != nil).To(Equal(tc.err))
			Expect(out).To(Equal(tc.expectedResult))
		}
	})
})
