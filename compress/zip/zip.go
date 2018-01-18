package zip

import (
	"os"
	"os/exec"

	"github.com/mtfelian/utils"
)

const (
	zipPath   = "zip"
	unzipPath = "unzip"
)

// Compress создаёт архив с путём и именем fileOut из файлов и папок, переданных в pathIn
func Compress(fileOut string, pathIn ...string) error {
	if utils.FileExists(fileOut) {
		if err := os.Remove(fileOut); err != nil {
			return err
		}
	}
	cmdParams := append([]string{"-j", fileOut}, pathIn...)
	_, err := exec.Command(zipPath, cmdParams...).CombinedOutput()
	return err
}

// Decompress извлекает из архива с путём и именем fileIn содержимое и помещает его в pathOut
func Decompress(fileIn string, pathOut string) error {
	cmdParams := []string{fileIn, "-d", pathOut}
	_, err := exec.Command(unzipPath, cmdParams...).CombinedOutput()
	return err
}
