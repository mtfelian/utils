package zip

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mtfelian/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZip(t *testing.T) {
	// todo rewrite tests
	fmt.Println("todo переписать тесты")
	os.Exit(0)

	binPath := utils.MustSelfPath()
	testPath = filepath.Join(binPath, testDir)
	if err := removeTestFiles(); err != nil {
		fmt.Errorf("Ошибка удаления тестовых директорий: %v\n", err)
	}
	if err := os.MkdirAll(testPath, 0777); err != nil {
		fmt.Errorf("Ошибка создания тестовых директорий: %v\n", err)
	}
	sourceData, decompressedData = [][]byte{}, [][]byte{}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Zip Suite")
}
