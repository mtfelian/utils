package gost

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mtfelian/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with Ginkgo", func() {
	AfterEach(func() { cleanup() })

	It("checks SignDER", func() {
		if utils.IsInVexor() {
			fmt.Println("No GOST OpenSSL in Vexor")
			return
		}

		createTestFile()

		testFileIn := filepath.Join(testPath, testFile)
		testFileOut := testFileIn + ".sgn"

		signer := NewGostSignerEncryptor(getSSLParams())
		Expect(signer.SignDER(testFileIn, testFileOut)).To(Succeed())
		Expect(utils.FileExists(testFileOut)).To(BeTrue())
		Expect(readData(&decryptedData)).To(Succeed())
	})

	It("encrypt", func() {
		if utils.IsInVexor() {
			fmt.Println("No GOST OpenSSL in Vexor")
			return
		}

		createTestFile()

		testFileIn := filepath.Join(testPath, testFile)
		testFileOut := testFileIn + ".enc"

		signer := NewGostSignerEncryptor(getSSLParams())
		Expect(signer.Encrypt(testFileIn, testFileOut)).To(Succeed())
		Expect(utils.FileExists(testFileOut)).To(BeTrue())
		Expect(readData(&decryptedData)).To(Succeed())
	})
})

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

func removeTestFiles() error { return os.RemoveAll(testPath) }

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

func writeTestCerts() {
	Expect(ioutil.WriteFile(filepath.Join(testPath, "private.pem"), []byte(testPrivateKey), 0660)).To(Succeed())
	Expect(ioutil.WriteFile(filepath.Join(testPath, "cert.cer"), []byte(testCertificate), 0660)).To(Succeed())
}

func createTestFile() {
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
	Expect(ioutil.WriteFile(file1Path, []byte("file1data"), 0660)).To(Succeed())
	testFile = file1Name

	Expect(readData(&sourceData)).To(Succeed())
	writeTestCerts()
}

func getSSLParams() SSLParams {
	return SSLParams{
		OpenSSLPath:         "/gost-ssl/bin/openssl",
		OurCertFilePath:     filepath.Join(testPath, "cert.cer"),
		ForeignCertFilePath: filepath.Join(testPath, "cert.cer"),
		OurPrivateKey:       filepath.Join(testPath, "private.pem"),
	}
}

func cleanup() { removeTestFiles() }
