package gost

import (
	"fmt"
	"github.com/mihteh/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testDir = "test"

const testPrivateKey = `-----BEGIN PRIVATE KEY-----
MEUCAQAwHAYGKoUDAgITMBIGByqFAwICIwEGByqFAwICHgEEIgIgTAJW4VqnYtYP
sPL4CN3b08ZolXU7iYN0CAWPGH7GKvY=
-----END PRIVATE KEY-----`

const testCertificate = `-----BEGIN CERTIFICATE-----
MIICtzCCAmQCCQC7UQoWC1zZ6zAKBgYqhQMCAgMFADCB4TELMAkGA1UEBhMCUlUx
LDAqBgNVBAgMI9Ca0YDQsNGB0L3QvtC00LDRgNGB0LrQuNC5INC60YDQsNC5MRsw
GQYDVQQHDBLQmtGA0LDRgdC90L7QtNCw0YAxKjAoBgNVBAoMIdCe0J7QniAi0KDQ
vtCz0LAg0Lgg0LrQvtC/0YvRgtCwIjEYMBYGA1UECwwP0JjQoi3QvtGC0LTQtdC7
MSowKAYDVQQDDCHQntCe0J4gItCg0L7Qs9CwINC4INC60L7Qv9GL0YLQsCIxFTAT
BgkqhkiG9w0BCQEWBmFAYS5ydTAeFw0xNjA4MDUwODM4MTFaFw0xOTAxMjIwODM4
MTFaMIHhMQswCQYDVQQGEwJSVTEsMCoGA1UECAwj0JrRgNCw0YHQvdC+0LTQsNGA
0YHQutC40Lkg0LrRgNCw0LkxGzAZBgNVBAcMEtCa0YDQsNGB0L3QvtC00LDRgDEq
MCgGA1UECgwh0J7QntCeICLQoNC+0LPQsCDQuCDQutC+0L/Ri9GC0LAiMRgwFgYD
VQQLDA/QmNCiLdC+0YLQtNC10LsxKjAoBgNVBAMMIdCe0J7QniAi0KDQvtCz0LAg
0Lgg0LrQvtC/0YvRgtCwIjEVMBMGCSqGSIb3DQEJARYGYUBhLnJ1MGMwHAYGKoUD
AgITMBIGByqFAwICIwEGByqFAwICHgEDQwAEQNz1pBRgNdbt0/EmcAHtKX83YYS7
JArNQBDAk7QMkGr2XNqpe8FFCqvH6aIeMPwcJ/CmpOR60Rugf3N4FJ3VA+AwCgYG
KoUDAgIDBQADQQBXjYwQLuQqhdS7yhY0gh38behzPIdWUQaPZIu/+BYvZF8szXdu
ID4lpcoZxRwQ37jX+suvd6koFC6V00gEnRCo
-----END CERTIFICATE-----`

var (
	testFile string
	testPath string

	sourceData    []byte
	decryptedData []byte
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}


func removeTestFiles() error {
	return os.RemoveAll(testPath)
}

// readData читает данные из получившегося в результате проверяемых манипуляций файла в тестовые переменные
func readData(into *[]byte) error {
	*into = []byte{}
	p := filepath.Join(testPath, testFile)
	isDir, err := utils.IsDir(p)
	if err != nil {
		return fmt.Errorf("Ошибка IsDir() файла %s: %v", p, err)
	}
	if isDir {
		return fmt.Errorf("Тестовый файл %s оказался директорией: %v", p, err)
	}
	fileData, err := ioutil.ReadFile(p)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла %s: %v", p, err)
	}
	*into = fileData
	return nil
}

// writeTestCerts записывает тестовый приватный ключ и сертификат
func writeTestCerts(t *testing.T) {
	if err := ioutil.WriteFile(filepath.Join(testPath, "private.pem"), []byte(testPrivateKey), 0660); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(testPath, "cert.cer"), []byte(testCertificate), 0660); err != nil {
		t.Fatal(err)
	}
}

// createTestFile создаёт тестовый файл
func createTestFile(t *testing.T) {
	binPath, err := utils.GetSelfPath()
	if err != nil {
		fmt.Errorf("Ошибка получения пути к приложению: %v\n", err)
	}
	testPath = filepath.Join(binPath, testDir)
	if err := removeTestFiles(); err != nil {
		fmt.Errorf("Ошибка удаления тестовых директорий: %v\n", err)
	}
	if err := os.MkdirAll(testPath, 0777); err != nil {
		fmt.Errorf("Ошибка создания тестовых директорий: %v\n", err)
	}
	sourceData = []byte{}
	decryptedData = []byte{}

	testFile = ""

	file1Name := "file1"
	file1Path := filepath.Join(testPath, file1Name)
	if err := ioutil.WriteFile(file1Path, []byte("file1data"), 0660); err != nil {
		t.Fatal(err)
	}
	testFile = file1Name

	if err := readData(&sourceData); err != nil {
		t.Fatal(err)
	}

	writeTestCerts(t)
}

// getSSLParams возвращает параметры SSL для теста
func getSSLParams() SSLParams {
	return SSLParams{
		OpenSSLPath:         "/gost-ssl/bin/openssl",
		OurCertFilePath:     filepath.Join(testPath, "cert.cer"),
		ForeignCertFilePath: filepath.Join(testPath, "cert.cer"),
		OurPrivateKey:       filepath.Join(testPath, "private.pem"),
	}
}

// cleanup сбрасывает тестовые переменные и удаляет тестовые файлы и директории
func cleanup() {
	removeTestFiles()
}

func TestSignDER(t *testing.T) {
	if utils.IsInVexor() {
		fmt.Println("No GOST OpenSSL in Vexor")
		return
	}

	defer cleanup()
	createTestFile(t)

	testFileIn := filepath.Join(testPath, testFile)
	testFileOut := testFileIn + ".sgn"

	signer := NewGostSignerEncryptor(getSSLParams())
	if err := signer.SignDER(testFileIn, testFileOut); err != nil {
		t.Fatalf("Ошибка SignDER(): %v", err)
	}

	// проверяем только что подписанный файл существует и его получается прочитать
	if !utils.FileExists(testFileOut) {
		t.Errorf("Файл %s должен быть распакован, но не найден.\n", testFileOut)
	}
	if err := readData(&decryptedData); err != nil {
		t.Error(err)
	}

	// todo хорошо бы проверить снятие подписи
}

func TestEncrypt(t *testing.T) {
	if utils.IsInVexor() {
		fmt.Println("No GOST OpenSSL in Vexor")
		return
	}

	defer cleanup()
	createTestFile(t)

	testFileIn := filepath.Join(testPath, testFile)
	testFileOut := testFileIn + ".enc"

	signer := NewGostSignerEncryptor(getSSLParams())
	if err := signer.Encrypt(testFileIn, testFileOut); err != nil {
		t.Fatalf("Ошибка SignDER(): %v", err)
	}

	// проверяем только что подписанный файл существует и его получается прочитать
	if !utils.FileExists(testFileOut) {
		t.Errorf("Файл %s должен быть распакован, но не найден.\n", testFileOut)
	}
	if err := readData(&decryptedData); err != nil {
		t.Error(err)
	}

	// todo хорошо бы проверить расшифровку
}
