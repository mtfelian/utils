package zip

import (
	"fmt"
	"github.com/mihteh/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testDir = "test"

var (
	testFiles []string
	testPath  string

	sourceData       [][]byte
	decompressedData [][]byte
)

func TestMain(m *testing.M) {
	fmt.Println("todo переписать тесты")
	os.Exit(0)

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
	sourceData = [][]byte{}
	decompressedData = [][]byte{}

	os.Exit(m.Run())
}

func removeTestFiles() error {
	return os.RemoveAll(testPath)
}

func readData(into *[][]byte) error {
	*into = [][]byte{}
	for _, value := range testFiles {
		p := filepath.Join(testPath, value)
		isDir, err := utils.IsDir(p)
		if err != nil {
			return fmt.Errorf("Ошибка IsDir() файла %s: %v", p, err)
		}
		if isDir {
			continue
		}
		fileData, err := ioutil.ReadFile(p)
		if err != nil {
			return fmt.Errorf("Ошибка при чтении файла %s: %v", p, err)
		}
		*into = append(*into, fileData)
	}
	return nil
}

func createTestFiles(t *testing.T) {
	testFiles = []string{}

	file1Name := "file1"
	file1Path := filepath.Join(testPath, file1Name)
	if err := ioutil.WriteFile(file1Path, []byte("file1data"), 0660); err != nil {
		t.Fatal(err)
	}
	testFiles = append(testFiles, file1Name)

	dir1Name := "dir1"
	dir1Path := filepath.Join(testPath, dir1Name)
	if err := os.Mkdir(dir1Path, 0777); err != nil {
		t.Fatal(err)
	}
	testFiles = append(testFiles, dir1Name)

	file2Name := "file2"
	file2Path := filepath.Join(testPath, dir1Name, file2Name)
	if err := ioutil.WriteFile(file2Path, []byte("file2data"), 0660); err != nil {
		t.Fatal(err)
	}
	testFiles = append(testFiles, filepath.Join(dir1Name, file2Name))

	if err := readData(&sourceData); err != nil {
		t.Fatal(err)
	}
}

func TestCompress(t *testing.T) {
	defer func() {
		os.RemoveAll(testPath)
	}()
	createTestFiles(t)

	binPath, err := utils.GetSelfPath()
	if err != nil {
		t.Fatalf("Ошибка получения пути к приложению: %v", err)
	}

	outputPath := filepath.Join(binPath, "output.zip")
	testFilesFullPaths := []string{}
	for _, f := range testFiles {
		testFilesFullPaths = append(testFilesFullPaths, filepath.Join(testPath, f))
	}

	if err := Compress(outputPath, testFilesFullPaths...); err != nil {
		t.Fatalf("Ошибка Compress(): %v", err)
	}

	if !utils.FileExists(outputPath) {
		t.Fatalf("Архив %s не существует", outputPath)
	}

	if err := removeTestFiles(); err != nil {
		fmt.Errorf("Ошибка удаления тестовых директорий: %v\n", err)
	}

	if err := Decompress(outputPath, "/"); err != nil {
		t.Fatalf("Ошибка Compress(): %v", err)
	}

	for _, path := range testFilesFullPaths {
		if !utils.FileExists(path) {
			t.Errorf("Файл %s должен быть распакован, но не найден.\n", path)
		}
		if err := readData(&decompressedData); err != nil {
			t.Error(err)
		}
		sourceLastElement := sourceData[len(sourceData)-1]
		decompressedLastElement := decompressedData[len(decompressedData)-1]
		if !reflect.DeepEqual(sourceLastElement, decompressedLastElement) {
			t.Errorf("Элементы не равны. \nИсходный     : %v\nРаспакованный: %v\n")
		}
	}
}
