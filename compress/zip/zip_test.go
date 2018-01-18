package zip

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
	BeforeEach(func() { createTestFiles() })
	AfterEach(func() { os.RemoveAll(testPath) })

	It("checks Compress/Decompress", func() {
		binPath, err := utils.GetSelfPath()
		Expect(err).NotTo(HaveOccurred())

		outputPath := filepath.Join(binPath, "output.zip")
		testFilesFullPaths := []string{}
		for _, f := range testFiles {
			testFilesFullPaths = append(testFilesFullPaths, filepath.Join(testPath, f))
		}

		Expect(Compress(outputPath, testFilesFullPaths...)).To(Succeed())
		Expect(utils.FileExists(outputPath)).To(BeTrue())
		Expect(removeTestFiles()).To(Succeed())
		Expect(Decompress(outputPath, "/")).To(Succeed())

		for _, path := range testFilesFullPaths {
			By(fmt.Sprintf("testing case %s", path))
			Expect(utils.FileExists(path)).To(BeTrue())
			Expect(readData(&decompressedData)).To(Succeed())
			sourceLastElement := sourceData[len(sourceData)-1]
			decompressedLastElement := decompressedData[len(decompressedData)-1]
			Expect(decompressedLastElement).To(Equal(sourceLastElement))
		}
	})
})

const testDir = "test"

var (
	testFiles []string
	testPath  string

	sourceData       [][]byte
	decompressedData [][]byte
)

func removeTestFiles() error { return os.RemoveAll(testPath) }

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

func createTestFiles() {
	testFiles = []string{}

	file1Name := "file1"
	file1Path := filepath.Join(testPath, file1Name)
	Expect(ioutil.WriteFile(file1Path, []byte("file1data"), 0660)).To(Succeed())
	testFiles = append(testFiles, file1Name)

	dir1Name := "dir1"
	dir1Path := filepath.Join(testPath, dir1Name)
	Expect(os.Mkdir(dir1Path, 0777)).To(Succeed())
	testFiles = append(testFiles, dir1Name)

	file2Name := "file2"
	file2Path := filepath.Join(testPath, dir1Name, file2Name)
	Expect(ioutil.WriteFile(file2Path, []byte("file2data"), 0660)).To(Succeed())
	testFiles = append(testFiles, filepath.Join(dir1Name, file2Name))

	Expect(readData(&sourceData)).To(Succeed())
}
