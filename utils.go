package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mtfelian/validation"
	"golang.org/x/text/encoding/charmap"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// IsInVexor возвращает true если выполнение происходит в среде Vexor, иначе false
func IsInVexor() bool {
	return os.Getenv("CI_NAME") == "VEXOR"
}

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

	randInt := func(min int, max int) int { return min + rand.Intn(max-min) }

	symbols := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for i := 0; i < n; i++ {
		b[i] = symbols[byte(randInt(0, len(symbols)))]
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

// trimSnils удаляет из СНИЛС всё кроме цифр
func trimSnils(snils string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(snils, "")
}

// CheckSnils проверяет СНИЛС на валидность путём вычисления его контрольной суммы
func CheckSnils(snils string) (bool, error) {
	const minimumSnilsCanValidate = 1001998

	s := trimSnils(snils)
	pattern := regexp.MustCompile(`^\d{11}$`)
	if !pattern.MatchString(s) {
		return false, nil
	}

	digits := strings.Split(s, "") // все цифры СНИЛСа
	numberDigits := digits[0:9]    // цифры номера
	numUint, err := strconv.ParseUint(strings.Join(numberDigits, ""), 10, 32)
	if err != nil {
		return false, err
	}

	// номер 0 не валидный
	if numUint == 0 {
		return false, nil
	}

	// считаем валидными те номера, для которых не считается контрольная сумма
	if numUint < minimumSnilsCanValidate {
		return true, nil
	}

	checkSumDigits := digits[9:11] // цифры контрольной суммы
	checkSumUint, err := strconv.ParseUint(strings.Join(checkSumDigits, ""), 10, 32)
	if err != nil {
		return false, err
	}

	sum, i := uint(0), uint(9)
	for _, digit := range numberDigits {
		digitAsUint, err := strconv.ParseUint(digit, 10, 32)
		if err != nil {
			return false, nil
		}
		sum += uint(digitAsUint) * i
		i--
	}

	expectedCheckSum := sum % 101
	if expectedCheckSum == 100 {
		expectedCheckSum = 0
	}

	return expectedCheckSum == uint(checkSumUint), nil
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

// EncodeToWindows1251 перекодирует срез байт b из стандартной Go кодировки UTF-8
// в кодировку Windows-1251
func EncodeToWindows1251(b []byte) ([]byte, error) {
	enc := charmap.Windows1251.NewEncoder()
	out, err := enc.Bytes(b)
	if err != nil {
		return []byte(""), err
	}
	return out, nil
}

// CountElementsOnPage возвращает количество элементов на заданной странице page с размером pageSize
// если всего элементов elementsTotal. Если pageSize равно 0, считается что оно равно elementsTotal
func CountElementsOnPage(elementsTotal uint, page uint, pageSize uint) uint {
	if pageSize < 1 {
		pageSize = elementsTotal
	}
	if page < 1 {
		page = 1
	}

	pages := CountPages(elementsTotal, pageSize)
	elementsOnLastPage := elementsTotal % pageSize

	if page > pages {
		return 0
	}

	if elementsOnLastPage == 0 {
		elementsOnLastPage = pageSize
	}

	if page == pages {
		return elementsOnLastPage
	}

	return pageSize
}

// CountPages возвращает количство страниц размера pageSize если всего элементов elementsTotal.
// Если pageSize равно 0, возвращает 1
func CountPages(elementsTotal uint, pageSize uint) uint {
	if elementsTotal == 0 {
		return 0
	}

	if pageSize == 0 {
		return 1
	}

	countPages := elementsTotal / pageSize
	if elementsTotal%pageSize != 0 {
		countPages++
	}
	return countPages
}

// StringSlice это срез строк
// реализует интерфейс Stringer
type StringSlice []string

// String возвращает строку, содержащую значения среза строк,
// где элементы разделены переносами строки
func (ss StringSlice) String() string {
	result := ""
	for _, element := range ss {
		result += fmt.Sprintf("%s\n", element)
	}
	return result
}

// IsNil возвращает true если объект является nil или содержит значение, эквивалентное нулевому
// иначе возвращает false
func IsNil(obj interface{}) bool {
	return !(validation.Required{}).IsSatisfied(obj)
}

// PString возвращает указатель на строку s или nil, если строка пустая
func PString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PUint возвращает указатель на uint i
func PUint(i uint) *uint {
	return &i
}

// PInt возвращает указатель на int i
func PInt(i int) *int {
	return &i
}

// GetIPAddress пытается получить IP адрес из заголовков HTTP
// возвращает соотв-ю строку, или "0.0.0.0"
func GetIPAddress(request *http.Request) string {
	regexpIP4 := regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.` +
		`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.` +
		`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.` +
		`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

	ipKeys := []string{"HTTP_CLIENT_IP",
		"HTTP_X_FORWARDED_FOR",
		"HTTP_X_FORWARDED",
		"HTTP_X_CLUSTER_CLIENT_IP",
		"HTTP_FORWARDED_FOR",
		"HTTP_FORWARDED",
		"REMOTE_ADDR",
	}

	for _, headerKey := range ipKeys {
		headerValue := request.Header.Get(headerKey)
		if headerValue != "" {
			headerValueParts := strings.Split(headerValue, ",")
			for _, headerValuePart := range headerValueParts {
				ip := strings.TrimSpace(headerValuePart)
				if regexpIP4.MatchString(ip) {
					return ip
				}
			}
		}
	}

	remoteAddr := request.Header.Get("REMOTE_ADDR")
	if remoteAddr != "" {
		return remoteAddr
	}
	return "0.0.0.0"
}

// CallerFuncName возвращает имя функции, вызвавшей функцию, из которой была вызвана CallerFuncName()
func CallerFuncName() (string, error) {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "", errors.New("Error after runtime.Callers(), n == 0")
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "", errors.New("Error after runtime.FuncForPC(): fun == nil")
	}

	return fun.Name(), nil
}

// CallerFuncNameString returns a string from
func CallerFuncNameString() string {
	funcName, err := CallerFuncName()
	if err != nil {
		return ""
	}
	return funcName
}

// RemoveDuplicates удаляет повторные значения из среза s
func RemoveDuplicates(s *[]uint) {
	found := make(map[uint]bool)
	j := 0
	for i, element := range *s {
		if !found[element] {
			found[element] = true
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

// UintSlice attaches the methods of sort.Interface to []uint, sorting in increasing order.
type UintSlice []uint

func (p UintSlice) Len() int           { return len(p) }
func (p UintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p UintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p UintSlice) Sort()              { sort.Sort(p) }

// SortUints sorts a slice of uints in increasing order.
func SortUints(a []uint) { sort.Sort(UintSlice(a)) }

// UintsAreSorted tests whether a slice of uints is sorted in increasing order.
func UintsAreSorted(a []uint) bool { return sort.IsSorted(UintSlice(a)) }

// FileUploadRequest это параметры для POST запроса с файлом
type FileUploadRequest struct {
	Uri      string            // uri to send request
	Params   map[string]string // additional parameters or nil, would be written into request fields
	Key      string            // key of multipart field
	Data     []byte            // file data
	FileName string            // file name
}

// NewFileUploadRequest creates a new file upload HTTP request with optional extra params
func NewFileUploadRequest(req FileUploadRequest) (*http.Request, error) {
	body := bytes.NewBuffer([]byte{})

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(req.Key, req.FileName)
	if err != nil {
		return nil, err
	}
	part.Write(req.Data)

	for key, val := range req.Params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", req.Uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

// ParseCookieString разбирает строку с куками и возвращает результат в виде среза указателей на куки
func ParseCookieString(cookie string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", cookie)
	request := http.Request{Header: header}
	return request.Cookies()
}
