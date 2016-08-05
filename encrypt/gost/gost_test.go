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

var (
	testFile string
	testPath string

	sourceData    []byte
	decryptedData []byte
)

func TestMain(m *testing.M) {
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

	os.Exit(m.Run())
}

func removeTestFiles() error {
	return os.RemoveAll(testPath)
}

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

func createTestFile(t *testing.T) {
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
}

func getSSLParams() SSLParams {
	return SSLParams{
		OpenSSLPath:         "/gost-ssl/bin/openssl",
		OurCertFilePath:     "/gfk/certs/kakunin/kakunin.cer",
		ForeignCertFilePath: "/gfk/certs/equifax1617/prod/Prod_Equifax_2016-2017.cer",
		OurPrivateKey:       "/gfk/private.pem",
	}
}

func TestSignDER(t *testing.T) {
	defer func() {
		os.RemoveAll(testPath)
	}()
	createTestFile(t)

	testFileIn := filepath.Join(testPath, testFile)
	testFileOut := testFileIn + ".sgn"

	signer := NewGostSigner(getSSLParams())
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

}
